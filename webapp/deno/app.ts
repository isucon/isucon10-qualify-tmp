import { Application, Router, helpers } from "https://deno.land/x/oak/mod.ts";
import { Client } from "https://deno.land/x/mysql/mod.ts";
import { sprintf } from "https://deno.land/std/fmt/printf.ts";
import { organ } from "https://raw.githubusercontent.com/denjucks/organ/master/mod.ts";
import { camelCase } from "https://deno.land/x/case/mod.ts";
import { parse } from "https://deno.land/std/encoding/csv.ts";
import { multiParser } from "https://deno.land/x/multiparser/mod.ts";

const currentEnv = Deno.env.toObject();
const decoder = new TextDecoder();

const PORT = currentEnv.PORT ?? 1323;
const LIMIT = 20;
const NAZOTTE_LIMIT = 50;
const dbinfo = {
  hostname: currentEnv.MYSQL_HOST ?? "127.0.0.1",
  port: parseInt(currentEnv.MYSQL_PORT ?? 3306),
  username: currentEnv.MYSQL_USER ?? "isucon",
  password: currentEnv.MYSQL_PASS ?? "isucon",
  db: currentEnv.MYSQL_DBNAME ?? "isuumo",
  poolSize: 10,
  debug: true,
};

const chairSearchConditionJson = await Deno.readFile(
  "../fixture/chair_condition.json",
);
const chairSearchCondition = JSON.parse(
  decoder.decode(chairSearchConditionJson),
);

const estateSearchConditionJson = await Deno.readFile(
  "../fixture/estate_condition.json",
);
const estateSearchCondition = JSON.parse(
  decoder.decode(estateSearchConditionJson),
);

const db = await new Client().connect(dbinfo);

const camelcaseKeys = (obj: any) =>
  Object.fromEntries(Object.entries(obj).map(([k, v]) => [camelCase(k), v]));

const router = new Router();

router.post("/initialize", async (ctx) => {
  const dbdir = "../mysql/db";
  const dbfiles = [
    "0_Schema.sql",
    "1_DummyEstateData.sql",
    "2_DummyChairData.sql",
  ];
  const execfiles = dbfiles.map((file) => `${dbdir}/${file}`);
  for (const execfile of execfiles) {
    const p = Deno.run({
      cmd: [
        "bash",
        "-c",
        `mysql -h ${dbinfo.hostname} -u ${dbinfo.username} -p${dbinfo.password} -P ${dbinfo.port} ${dbinfo.db} < ${execfile}`,
      ],
    });
    const status = await p.status();
    if (!status.success) {
      const output = await p.output();
      throw new Error("Deno run is failed " + output);
    }
  }
  ctx.response.body = {
    language: "deno",
  };
});

router.get("/api/estate/low_priced", async (ctx) => {
  const es = await db.query(
    "SELECT * FROM estate ORDER BY rent ASC, id ASC LIMIT ?",
    [LIMIT],
  );
  ctx.response.body = { estates: es.map(camelcaseKeys) };
});

router.get("/api/chair/low_priced", async (ctx) => {
  try {
    const cs = await db.query(
      "SELECT * FROM chair WHERE stock > 0 ORDER BY price ASC, id ASC LIMIT ?",
      [LIMIT],
    );
    ctx.response.body = { chairs: cs.map(camelcaseKeys) };
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.get("/api/chair/search", async (ctx, next) => {
  const searchQueries = [];
  const queryParams = [];
  const {
    priceRangeId,
    heightRangeId,
    widthRangeId,
    depthRangeId,
    kind,
    color,
    features,
    page,
    perPage,
  } = helpers.getQuery(ctx);

  if (!!priceRangeId) {
    const chairPrice = chairSearchCondition["price"].ranges[priceRangeId];
    if (chairPrice == null) {
      ctx.response.status = 400;
      ctx.response.body = "priceRangeID invalid";
      return;
    }

    if (chairPrice.min !== -1) {
      searchQueries.push("price >= ? ");
      queryParams.push(chairPrice.min);
    }

    if (chairPrice.max !== -1) {
      searchQueries.push("price < ? ");
      queryParams.push(chairPrice.max);
    }
  }

  if (!!heightRangeId) {
    const chairHeight = chairSearchCondition["height"].ranges[heightRangeId];
    if (chairHeight == null) {
      ctx.response.status = 400;
      ctx.response.body = "heightRangeId invalid";
      return;
    }

    if (chairHeight.min !== -1) {
      searchQueries.push("height >= ? ");
      queryParams.push(chairHeight.min);
    }

    if (chairHeight.max !== -1) {
      searchQueries.push("height < ? ");
      queryParams.push(chairHeight.max);
    }
  }

  if (!!widthRangeId) {
    const chairWidth = chairSearchCondition["width"].ranges[widthRangeId];
    if (chairWidth == null) {
      ctx.response.status = 400;
      ctx.response.body = "widthRangeId invalid";
      return;
    }

    if (chairWidth.min !== -1) {
      searchQueries.push("width >= ? ");
      queryParams.push(chairWidth.min);
    }

    if (chairWidth.max !== -1) {
      searchQueries.push("width < ? ");
      queryParams.push(chairWidth.max);
    }
  }

  if (!!depthRangeId) {
    const chairDepth = chairSearchCondition["depth"].ranges[depthRangeId];
    if (chairDepth == null) {
      ctx.response.status = 400;
      ctx.response.body = "depthRangeId invalid";
      return;
    }

    if (chairDepth.min !== -1) {
      searchQueries.push("depth >= ? ");
      queryParams.push(chairDepth.min);
    }

    if (chairDepth.max !== -1) {
      searchQueries.push("depth < ? ");
      queryParams.push(chairDepth.max);
    }
  }

  if (!!kind) {
    searchQueries.push("kind = ? ");
    queryParams.push(kind);
  }

  if (!!color) {
    searchQueries.push("color = ? ");
    queryParams.push(color);
  }

  if (!!features) {
    const featureConditions = features.split(",");
    for (const featureCondition of featureConditions) {
      searchQueries.push("features LIKE ?");
      queryParams.push(`%${featureCondition}%`);
    }
  }

  if (searchQueries.length === 0) {
    ctx.response.status = 400;
    ctx.response.body = "Search condition not found";
    return;
  }

  searchQueries.push("stock > 0");
  const pageNum = parseInt(page, 10);
  const perPageNum = parseInt(perPage, 10);

  if (!page || Number.isNaN(pageNum)) {
    ctx.response.status = 400;
    ctx.response.body = `page condition invalid ${page}`;
    return;
  }

  if (!perPage || Number.isNaN(perPageNum)) {
    ctx.response.status = 400;
    ctx.response.body = `perPage condition invalid ${perPage}`;
    return;
  }

  const sqlprefix = "SELECT * FROM chair WHERE ";
  const searchCondition = searchQueries.join(" AND ");
  const limitOffset = " ORDER BY popularity DESC, id ASC LIMIT ? OFFSET ?";
  const countprefix = "SELECT COUNT(*) as count FROM chair WHERE ";

  try {
    const [{ count }] = await db.query(
      `${countprefix}${searchCondition}`,
      queryParams,
    );
    queryParams.push(perPageNum, perPageNum * pageNum);
    const cs = await db.query(
      `${sqlprefix}${searchCondition}${limitOffset}`,
      queryParams,
    );
    const chairs = cs.map(camelcaseKeys);
    ctx.response.body = {
      count,
      chairs,
    };
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.get("/api/chair/search/condition", (ctx) => {
  ctx.response.body = chairSearchCondition;
});

router.get("/api/chair/:id", async (ctx) => {
  try {
    const id = ctx.params.id;
    const [chair] = await db.query("SELECT * FROM chair WHERE id = ?", [id]);
    if (chair == null || chair.stock <= 0) {
      ctx.response.status = 404;
      ctx.response.body = "Not Found";
      return;
    }
    ctx.response.body = camelcaseKeys(chair);
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.post("/api/chair/buy/:id", async (ctx) => {
  try {
    const id = ctx.params.id;
    await db.transaction(async (conn) => {
      const result = await conn.execute(
        "SELECT * FROM chair WHERE id = ? AND stock > 0 FOR UPDATE",
        [id],
      );
      if (result.rows?.[0] == null) {
        ctx.response.status = 404;
        ctx.response.body = "Not Found";
        await conn.execute("ROLLBACK");
        return;
      }
      const chair = result.rows[0];
      await conn.execute(
        "UPDATE chair SET stock = ? WHERE id = ?",
        [chair.stock - 1, id],
      );
    });

    ctx.response.body = { ok: true };
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.get("/api/estate/search", async (ctx) => {
  const searchQueries = [];
  const queryParams = [];
  const {
    doorHeightRangeId,
    doorWidthRangeId,
    rentRangeId,
    features,
    page,
    perPage,
  } = helpers.getQuery(ctx);

  if (!!doorHeightRangeId) {
    const doorHeight =
      estateSearchCondition["doorHeight"].ranges[doorHeightRangeId];
    if (doorHeight == null) {
      ctx.response.status = 400;
      ctx.response.body = "doorHeightRangeId invalid";
      return;
    }

    if (doorHeight.min !== -1) {
      searchQueries.push("door_height >= ? ");
      queryParams.push(doorHeight.min);
    }

    if (doorHeight.max !== -1) {
      searchQueries.push("door_height < ? ");
      queryParams.push(doorHeight.max);
    }
  }

  if (!!doorWidthRangeId) {
    const doorWidth =
      estateSearchCondition["doorWidth"].ranges[doorWidthRangeId];
    if (doorWidth == null) {
      ctx.response.status = 400;
      ctx.response.body = "doorWidthRangeId invalid";
      return;
    }

    if (doorWidth.min !== -1) {
      searchQueries.push("door_width >= ? ");
      queryParams.push(doorWidth.min);
    }

    if (doorWidth.max !== -1) {
      searchQueries.push("door_width < ? ");
      queryParams.push(doorWidth.max);
    }
  }

  if (!!rentRangeId) {
    const rent = estateSearchCondition["rent"].ranges[rentRangeId];
    if (rent == null) {
      ctx.response.status = 400;
      ctx.response.body = "rentRangeId invalid";
      return;
    }

    if (rent.min !== -1) {
      searchQueries.push("rent >= ? ");
      queryParams.push(rent.min);
    }

    if (rent.max !== -1) {
      searchQueries.push("rent < ? ");
      queryParams.push(rent.max);
    }
  }

  if (!!features) {
    const featureConditions = features.split(",");
    for (const featureCondition of featureConditions) {
      searchQueries.push("features LIKE ?");
      queryParams.push(`%${featureCondition}%`);
    }
  }

  if (searchQueries.length === 0) {
    ctx.response.status = 400;
    ctx.response.body = "Search condition not found";
    return;
  }

  const pageNum = parseInt(page, 10);
  const perPageNum = parseInt(perPage, 10);

  if (!page || Number.isNaN(pageNum)) {
    ctx.response.status = 400;
    ctx.response.body = `page condition invalid ${page}`;
    return;
  }

  if (!perPage || Number.isNaN(perPageNum)) {
    ctx.response.status = 400;
    ctx.response.body = `perPage condition invalid ${perPage}`;
    return;
  }

  const sqlprefix = "SELECT * FROM estate WHERE ";
  const searchCondition = searchQueries.join(" AND ");
  const limitOffset = " ORDER BY popularity DESC, id ASC LIMIT ? OFFSET ?";
  const countprefix = "SELECT COUNT(*) as count FROM estate WHERE ";

  try {
    const [{ count }] = await db.query(
      `${countprefix}${searchCondition}`,
      queryParams,
    );
    queryParams.push(perPageNum, perPageNum * pageNum);
    const estates = await db.query(
      `${sqlprefix}${searchCondition}${limitOffset}`,
      queryParams,
    );
    ctx.response.body = {
      count,
      estates: estates.map(camelcaseKeys),
    };
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.get("/api/estate/search/condition", (ctx) => {
  ctx.response.body = estateSearchCondition;
});

router.post("/api/estate/req_doc/:id", async (ctx) => {
  const id = ctx.params.id;
  const [estate] = await db.query("SELECT * FROM estate WHERE id = ?", [id]);
  if (estate == null) {
    ctx.response.status = 404;
    ctx.response.body = "Not Found";
    return;
  }
  ctx.response.body = { ok: true };
});

router.post("/api/estate/nazotte", async (ctx) => {
  const result = ctx.request.body(); // content type automatically detected
  let coordinates;
  if (result.type === "json") {
    const val = await result.value; // an object of parsed JSON
    coordinates = val.coordinates;
  }
  const longitudes = coordinates.map((c: { longitude: number }) => c.longitude);
  const latitudes = coordinates.map((c: { latitude: number }) => c.latitude);
  const boundingbox = {
    topleft: {
      longitude: Math.min(...longitudes),
      latitude: Math.min(...latitudes),
    },
    bottomright: {
      longitude: Math.max(...longitudes),
      latitude: Math.max(...latitudes),
    },
  };

  try {
    const estates = await db.query(
      "SELECT * FROM estate WHERE latitude <= ? AND latitude >= ? AND longitude <= ? AND longitude >= ? ORDER BY popularity DESC, id ASC",
      [
        boundingbox.bottomright.latitude,
        boundingbox.topleft.latitude,
        boundingbox.bottomright.longitude,
        boundingbox.topleft.longitude,
      ],
    );

    const estatesInPolygon = [];
    for (const estate of estates) {
      const point = sprintf(
        "'POINT(%f %f)'",
        estate.latitude,
        estate.longitude,
      );
      const sql =
        "SELECT * FROM estate WHERE id = ? AND ST_Contains(ST_PolygonFromText(%s), ST_GeomFromText(%s))";
      const coordinatesToText = sprintf(
        "'POLYGON((%s))'",
        coordinates.map((coordinate: { latitude: number; longitude: number }) =>
          sprintf("%f %f", coordinate.latitude, coordinate.longitude)
        ).join(","),
      );
      const sqlstr = sprintf(sql, coordinatesToText, point);
      const [e] = await db.query(sqlstr, [estate.id]);
      if (e && Object.keys(e).length > 0) {
        estatesInPolygon.push(e);
      }
    }

    const results = {
      count: 0,
      estates: [] as Array<any>,
    };
    let i = 0;
    for (const estate of estatesInPolygon) {
      if (i >= NAZOTTE_LIMIT) {
        break;
      }
      // camelize
      results.estates.push(camelcaseKeys(estate));
      i++;
    }
    results.count = results.estates.length;
    ctx.response.body = results;
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.get("/api/estate/:id", async (ctx) => {
  try {
    const id = ctx.params.id;
    const [estate] = await db.query("SELECT * FROM estate WHERE id = ?", [id]);
    if (estate == null) {
      ctx.response.status = 404;
      ctx.response.body = "Estate Not Found";
      return;
    }
    ctx.response.body = camelcaseKeys(estate);
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.post("/api/estate/req_doc/:id", async (ctx) => {
  const id = ctx.params.id;
  const [estate] = await db.query("SELECT * FROM estate WHERE id = ?", [id]);
  if (estate == null) {
    ctx.response.status = 404;
    ctx.response.body = "Not Found";
    return;
  }
  ctx.response.body = { ok: true };
});

router.get("/api/recommended_estate/:id", async (ctx) => {
  try {
    const id = ctx.params.id;
    const [chair] = await db.query("SELECT * FROM chair WHERE id = ?", [id]);
    const w = chair.width;
    const h = chair.height;
    const d = chair.depth;
    const es = await db.query(
      "SELECT * FROM estate where (door_width >= ? AND door_height>= ?) OR (door_width >= ? AND door_height>= ?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) OR (door_width >= ? AND door_height>=?) ORDER BY popularity DESC, id ASC LIMIT ?",
      [
        w,
        h,
        w,
        d,
        h,
        w,
        h,
        d,
        d,
        w,
        d,
        h,
        LIMIT,
      ],
    );
    ctx.response.body = { estates: es.map(camelcaseKeys) };
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.post("/api/chair", async (ctx) => {
  try {
    const form = await multiParser(ctx.request.serverRequest);
    if (!form || !form.chairs) {
      ctx.response.status = 400;
      ctx.response.body = "Bad Request";
      return;
    }
    const content = decoder.decode((form.chairs as any).content);
    const csv = await parse(content);
    await db.transaction(async (conn) => {
      for (let i = 0; i < csv.length; i++) {
        const items = csv[i] as any;
        await conn.execute(
          "INSERT INTO chair(id, name, description, thumbnail, price, height, width, depth, color, features, kind, popularity, stock) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)",
          items,
        );
      }
    });
    ctx.response.status = 201;
    ctx.response.body = { ok: true };
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});

router.post("/api/estate", async (ctx) => {
  try {
    const form = await multiParser(ctx.request.serverRequest);
    if (!form || !form.estates) {
      ctx.response.status = 400;
      ctx.response.body = "Bad Request";
      return;
    }
    const content = decoder.decode((form.estates as any).content);
    const csv = await parse(content);
    await db.transaction(async (conn) => {
      for (let i = 0; i < csv.length; i++) {
        const items = csv[i] as any;
        await conn.execute(
          "INSERT INTO estate(id, name, description, thumbnail, address, latitude, longitude, rent, door_height, door_width, features, popularity) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)",
          items,
        );
      }
    });
    ctx.response.status = 201;
    ctx.response.body = { ok: true };
  } catch (e) {
    ctx.response.status = 500;
    ctx.response.body = e.toString();
  }
});
const app = new Application();
app.use(organ());
app.use(router.routes());
app.use(router.allowedMethods());

console.log(`Listening ${PORT}`);
await app.listen({ port: +PORT });
await db.close();
