from os import getenv, path
import json
import subprocess
from io import StringIO
import csv
import flask
from werkzeug.exceptions import BadRequest, NotFound
import mysql.connector
from sqlalchemy.pool import QueuePool
from humps import camelize

LIMIT = 20
NAZOTTE_LIMIT = 50

chair_search_condition = json.load(open("../fixture/chair_condition.json", "r"))
estate_search_condition = json.load(open("../fixture/estate_condition.json", "r"))

app = flask.Flask(__name__)

mysql_connection_env = {
    "host": getenv("MYSQL_HOST", "127.0.0.1"),
    "port": getenv("MYSQL_PORT", 3306),
    "user": getenv("MYSQL_USER", "isucon"),
    "password": getenv("MYSQL_PASS", "isucon"),
    "database": getenv("MYSQL_DBNAME", "isuumo"),
}

cnxpool = QueuePool(lambda: mysql.connector.connect(**mysql_connection_env), pool_size=10)


def select_all(query, *args, dictionary=True):
    cnx = cnxpool.connect()
    try:
        cur = cnx.cursor(dictionary=dictionary)
        cur.execute(query, *args)
        return cur.fetchall()
    finally:
        cnx.close()


def select_row(*args, **kwargs):
    rows = select_all(*args, **kwargs)
    return rows[0] if len(rows) > 0 else None


@app.route("/initialize", methods=["POST"])
def post_initialize():
    sql_dir = "../mysql/db"
    sql_files = [
        "0_Schema.sql",
        "1_DummyEstateData.sql",
        "2_DummyChairData.sql",
    ]

    for sql_file in sql_files:
        command = f"mysql -h {mysql_connection_env['host']} -u {mysql_connection_env['user']} -p{mysql_connection_env['password']} -P {mysql_connection_env['port']} {mysql_connection_env['database']} < {path.join(sql_dir, sql_file)}"
        subprocess.run(["bash", "-c", command])

    return {"language": "python"}


@app.route("/api/estate/low_priced", methods=["GET"])
def get_estate_low_priced():
    rows = select_all("SELECT * FROM estate ORDER BY rent ASC, id ASC LIMIT %s", (LIMIT,))
    return {"estates": camelize(rows)}


@app.route("/api/chair/low_priced", methods=["GET"])
def get_chair_low_priced():
    rows = select_all("SELECT * FROM chair WHERE stock > 0 ORDER BY price ASC, id ASC LIMIT %s", (LIMIT,))
    return {"chairs": camelize(rows)}


@app.route("/api/chair/search", methods=["GET"])
def get_chair_search():
    args = flask.request.args

    conditions = []
    params = []

    if args.get("priceRangeId"):
        for _range in chair_search_condition["price"]["ranges"]:
            if _range["id"] == int(args.get("priceRangeId")):
                price = _range
                break
        else:
            raise BadRequest("priceRangeID invalid")
        if price["min"] != -1:
            conditions.append("price >= %s")
            params.append(price["min"])
        if price["max"] != -1:
            conditions.append("price < %s")
            params.append(price["max"])

    if args.get("heightRangeId"):
        for _range in chair_search_condition["height"]["ranges"]:
            if _range["id"] == int(args.get("heightRangeId")):
                height = _range
                break
        else:
            raise BadRequest("heightRangeId invalid")
        if height["min"] != -1:
            conditions.append("height >= %s")
            params.append(height["min"])
        if height["max"] != -1:
            conditions.append("height < %s")
            params.append(height["max"])

    if args.get("widthRangeId"):
        for _range in chair_search_condition["width"]["ranges"]:
            if _range["id"] == int(args.get("widthRangeId")):
                width = _range
                break
        else:
            raise BadRequest("widthRangeId invalid")
        if width["min"] != -1:
            conditions.append("width >= %s")
            params.append(width["min"])
        if width["max"] != -1:
            conditions.append("width < %s")
            params.append(width["max"])

    if args.get("depthRangeId"):
        for _range in chair_search_condition["depth"]["ranges"]:
            if _range["id"] == int(args.get("depthRangeId")):
                depth = _range
                break
        else:
            raise BadRequest("depthRangeId invalid")
        if depth["min"] != -1:
            conditions.append("depth >= %s")
            params.append(depth["min"])
        if depth["max"] != -1:
            conditions.append("depth < %s")
            params.append(depth["max"])

    if args.get("kind"):
        conditions.append("kind = %s")
        params.append(args.get("kind"))

    if args.get("color"):
        conditions.append("color = %s")
        params.append(args.get("color"))

    if args.get("features"):
        for feature_confition in args.get("features").split(","):
            conditions.append("features LIKE CONCAT('%', %s, '%')")
            params.append(feature_confition)

    if len(conditions) == 0:
        raise BadRequest("Search condition not found")

    conditions.append("stock > 0")

    try:
        page = int(args.get("page"))
    except (TypeError, ValueError):
        raise BadRequest("Invalid format page parameter")

    try:
        per_page = int(args.get("perPage"))
    except (TypeError, ValueError):
        raise BadRequest("Invalid format perPage parameter")

    search_condition = " AND ".join(conditions)

    query = f"SELECT COUNT(*) as count FROM chair WHERE {search_condition}"
    count = select_row(query, params)["count"]

    query = f"SELECT * FROM chair WHERE {search_condition} ORDER BY popularity DESC, id ASC LIMIT %s OFFSET %s"
    chairs = select_all(query, params + [per_page, per_page * page])

    return {"count": count, "chairs": camelize(chairs)}


@app.route("/api/chair/search/condition", methods=["GET"])
def get_chair_search_condition():
    return chair_search_condition


@app.route("/api/chair/<int:chair_id>", methods=["GET"])
def get_chair(chair_id):
    chair = select_row("SELECT * FROM chair WHERE id = %s", (chair_id,))
    if chair is None or chair["stock"] <= 0:
        raise NotFound()
    return camelize(chair)


@app.route("/api/chair/buy/<int:chair_id>", methods=["POST"])
def post_chair_buy(chair_id):
    cnx = cnxpool.connect()
    try:
        cnx.start_transaction()
        cur = cnx.cursor(dictionary=True)
        cur.execute("SELECT * FROM chair WHERE id = %s AND stock > 0 FOR UPDATE", (chair_id,))
        chair = cur.fetchone()
        if chair is None:
            raise NotFound()
        cur.execute("UPDATE chair SET stock = stock - 1 WHERE id = %s", (chair_id,))
        cnx.commit()
        return {"ok": True}
    except Exception as e:
        cnx.rollback()
        raise e
    finally:
        cnx.close()


@app.route("/api/estate/search", methods=["GET"])
def get_estate_search():
    args = flask.request.args

    conditions = []
    params = []

    if args.get("doorHeightRangeId"):
        for _range in estate_search_condition["doorHeight"]["ranges"]:
            if _range["id"] == int(args.get("doorHeightRangeId")):
                door_height = _range
                break
        else:
            raise BadRequest("doorHeightRangeId invalid")
        if door_height["min"] != -1:
            conditions.append("door_height >= %s")
            params.append(door_height["min"])
        if door_height["max"] != -1:
            conditions.append("door_height < %s")
            params.append(door_height["max"])

    if args.get("doorWidthRangeId"):
        for _range in estate_search_condition["doorWidth"]["ranges"]:
            if _range["id"] == int(args.get("doorWidthRangeId")):
                door_width = _range
                break
        else:
            raise BadRequest("doorWidthRangeId invalid")
        if door_width["min"] != -1:
            conditions.append("door_width >= %s")
            params.append(door_width["min"])
        if door_width["max"] != -1:
            conditions.append("door_width < %s")
            params.append(door_width["max"])

    if args.get("rentRangeId"):
        for _range in estate_search_condition["rent"]["ranges"]:
            if _range["id"] == int(args.get("rentRangeId")):
                rent = _range
                break
        else:
            raise BadRequest("rentRangeId invalid")
        if rent["min"] != -1:
            conditions.append("rent >= %s")
            params.append(rent["min"])
        if rent["max"] != -1:
            conditions.append("rent < %s")
            params.append(rent["max"])

    if args.get("features"):
        for feature_confition in args.get("features").split(","):
            conditions.append("features LIKE CONCAT('%', %s, '%')")
            params.append(feature_confition)

    if len(conditions) == 0:
        raise BadRequest("Search condition not found")

    try:
        page = int(args.get("page"))
    except (TypeError, ValueError):
        raise BadRequest("Invalid format page parameter")

    try:
        per_page = int(args.get("perPage"))
    except (TypeError, ValueError):
        raise BadRequest("Invalid format perPage parameter")

    search_condition = " AND ".join(conditions)

    query = f"SELECT COUNT(*) as count FROM estate WHERE {search_condition}"
    count = select_row(query, params)["count"]

    query = f"SELECT * FROM estate WHERE {search_condition} ORDER BY popularity DESC, id ASC LIMIT %s OFFSET %s"
    chairs = select_all(query, params + [per_page, per_page * page])

    return {"count": count, "estates": camelize(chairs)}


@app.route("/api/estate/search/condition", methods=["GET"])
def get_estate_search_condition():
    return estate_search_condition


@app.route("/api/estate/req_doc/<int:estate_id>", methods=["POST"])
def post_estate_req_doc(estate_id):
    estate = select_row("SELECT * FROM estate WHERE id = %s", (estate_id,))
    if estate is None:
        raise NotFound()
    return {"ok": True}


@app.route("/api/estate/nazotte", methods=["POST"])
def post_estate_nazotte():
    if "coordinates" not in flask.request.json:
        raise BadRequest()
    coordinates = flask.request.json["coordinates"]
    if len(coordinates) == 0:
        raise BadRequest()
    longitudes = [c["longitude"] for c in coordinates]
    latitudes = [c["latitude"] for c in coordinates]
    bounding_box = {
        "top_left_corner": {"longitude": min(longitudes), "latitude": min(latitudes)},
        "bottom_right_corner": {"longitude": max(longitudes), "latitude": max(latitudes)},
    }

    cnx = cnxpool.connect()
    try:
        cur = cnx.cursor(dictionary=True)
        cur.execute(
            (
                "SELECT * FROM estate"
                " WHERE latitude <= %s AND latitude >= %s AND longitude <= %s AND longitude >= %s"
                " ORDER BY popularity DESC, id ASC"
            ),
            (
                bounding_box["bottom_right_corner"]["latitude"],
                bounding_box["top_left_corner"]["latitude"],
                bounding_box["bottom_right_corner"]["longitude"],
                bounding_box["top_left_corner"]["longitude"],
            ),
        )
        estates = cur.fetchall()
        estates_in_polygon = []
        for estate in estates:
            query = "SELECT * FROM estate WHERE id = %s AND ST_Contains(ST_PolygonFromText(%s), ST_GeomFromText(%s))"
            polygon_text = (
                f"POLYGON(({','.join(['{} {}'.format(c['latitude'], c['longitude']) for c in coordinates])}))"
            )
            geom_text = f"POINT({estate['latitude']} {estate['longitude']})"
            cur.execute(query, (estate["id"], polygon_text, geom_text))
            if len(cur.fetchall()) > 0:
                estates_in_polygon.append(estate)
    finally:
        cnx.close()

    results = {"estates": []}
    for i, estate in enumerate(estates_in_polygon):
        if i >= NAZOTTE_LIMIT:
            break
        results["estates"].append(camelize(estate))
    results["count"] = len(results["estates"])
    return results


@app.route("/api/estate/<int:estate_id>", methods=["GET"])
def get_estate(estate_id):
    estate = select_row("SELECT * FROM estate WHERE id = %s", (estate_id,))
    if estate is None:
        raise NotFound()
    return camelize(estate)


@app.route("/api/recommended_estate/<int:chair_id>", methods=["GET"])
def get_recommended_estate(chair_id):
    chair = select_row("SELECT * FROM chair WHERE id = %s", (chair_id,))
    if chair is None:
        raise BadRequest(f"Invalid format searchRecommendedEstateWithChair id : {chair_id}")
    w, h, d = chair["width"], chair["height"], chair["depth"]
    query = (
        "SELECT * FROM estate"
        " WHERE (door_width >= %s AND door_height >= %s)"
        "    OR (door_width >= %s AND door_height >= %s)"
        "    OR (door_width >= %s AND door_height >= %s)"
        "    OR (door_width >= %s AND door_height >= %s)"
        "    OR (door_width >= %s AND door_height >= %s)"
        "    OR (door_width >= %s AND door_height >= %s)"
        " ORDER BY popularity DESC, id ASC"
        " LIMIT %s"
    )
    estates = select_all(query, (w, h, w, d, h, w, h, d, d, w, d, h, LIMIT))
    return {"estates": camelize(estates)}


@app.route("/api/chair", methods=["POST"])
def post_chair():
    if "chairs" not in flask.request.files:
        raise BadRequest()
    records = csv.reader(StringIO(flask.request.files["chairs"].read().decode()))
    cnx = cnxpool.connect()
    try:
        cnx.start_transaction()
        cur = cnx.cursor()
        for record in records:
            query = "INSERT INTO chair(id, name, description, thumbnail, price, height, width, depth, color, features, kind, popularity, stock) VALUES(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)"
            cur.execute(query, record)
        cnx.commit()
        return {"ok": True}, 201
    except Exception as e:
        cnx.rollback()
        raise e
    finally:
        cnx.close()


@app.route("/api/estate", methods=["POST"])
def post_estate():
    if "estates" not in flask.request.files:
        raise BadRequest()
    records = csv.reader(StringIO(flask.request.files["estates"].read().decode()))
    cnx = cnxpool.connect()
    try:
        cnx.start_transaction()
        cur = cnx.cursor()
        for record in records:
            query = "INSERT INTO estate(id, name, description, thumbnail, address, latitude, longitude, rent, door_height, door_width, features, popularity) VALUES(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)"
            cur.execute(query, record)
        cnx.commit()
        return {"ok": True}, 201
    except Exception as e:
        cnx.rollback()
        raise e
    finally:
        cnx.close()


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=getenv("SERVER_PORT", 1323), debug=True, threaded=True)
