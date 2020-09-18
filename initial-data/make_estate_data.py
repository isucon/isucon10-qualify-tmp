# coding:utf-8
import random
import json
import os
import glob
from faker import Faker
fake = Faker('ja_JP')
Faker.seed(19700101)
random.seed(19700101)

DESCRIPTION_LINES_FILE = "./description.txt"
OUTPUT_SQL_FILE = "./result/1_DummyEstateData.sql"
OUTPUT_TXT_FILE = "./result/estate_json.txt"
VERIFY_DRAFT_FILE = "./result/verify_draft_estate.txt"
OUTPUT_DRAFT_FILE = "./result/draft_data/estate/{index}.txt"
OUTPUT_FIXTURE_FILE = "./result/estate_condition.json"
ESTATE_IMAGE_ORIGIN_DIR = "./origin/estate"
ESTATE_IMAGE_PUBLIC_DIR = "../webapp/frontend/public/images/estate"
ESTATE_DUMMY_IMAGE_NUM = 1000
RECORD_COUNT = (10 ** 4) * 3 - 500
BULK_INSERT_COUNT = 500
DOOR_MIN_CENTIMETER = 30
DOOR_MAX_CENTIMETER = 200
DOOR_HEIGHT_RANGE_SEPARATORS = [80, 110, 150]
DOOR_WIDTH_RANGE_SEPARATORS = [80, 110, 150]
RENT_RANGE_SEPARATORS = [50000, 100000, 150000]
MIN_POPULARITY = 3000
MAX_POPULARITY = 1000000
VERIFY_DRAFT_COUNT = 500
DRAFT_COUNT_PER_FILE = 500
DRAFT_FILE_COUNT = 20

BUILDING_NAME_LIST = [
    "{name}ISUビルディング",
    "ISUアパート {name}",
    "ISU{name}レジデンス",
    "ISUガーデン {name}",
    "{name} ISUマンション",
    "{name} ISUビル"
]

ESTATE_FEATURE_LIST = [
    "最上階",
    "防犯カメラ",
    "ウォークインクローゼット",
    "ワンルーム",
    "ルーフバルコニー付",
    "エアコン付き",
    "駐輪場あり",
    "プロパンガス",
    "駐車場あり",
    "防音室",
    "追い焚き風呂",
    "オートロック",
    "即入居可",
    "IHコンロ",
    "敷地内駐車場",
    "トランクルーム",
    "角部屋",
    "カスタマイズ可",
    "DIY可",
    "ロフト",
    "シューズボックス",
    "インターネット無料",
    "地下室",
    "敷地内ゴミ置場",
    "管理人有り",
    "宅配ボックス",
    "ルームシェア可",
    "セキュリティ会社加入済",
    "メゾネット",
    "女性限定",
    "バイク置場あり",
    "エレベーター",
    "ペット相談可",
    "洗面所独立",
    "都市ガス",
    "浴室乾燥機",
    "インターネット接続可",
    "テレビ・通信",
    "専用庭",
    "システムキッチン",
    "高齢者歓迎",
    "ケーブルテレビ",
    "床下収納",
    "バス・トイレ別",
    "駐車場2台以上",
    "楽器相談可",
    "フローリング",
    "オール電化",
    "TVモニタ付きインタホン",
    "デザイナーズ物件"
]

ESTATE_IMAGE_HASH_LIST = [fake.sha256(
    raw_output=False) for _ in range(ESTATE_DUMMY_IMAGE_NUM)]


def generate_ranges_from_separator(separators):
    before = -1
    ranges = []

    for i, separator in enumerate(separators + [-1]):
        ranges.append({
            "id": i,
            "min": before,
            "max": separator
        })
        before = separator

    return ranges


def read_src_file_data(file_path):
    with open(file_path, mode='rb') as img:
        return img.read()


def dump_estate_to_json_str(estate):
    return json.dumps({
        "id": estate["id"],
        "thumbnail": estate["thumbnail"],
        "name": estate["name"],
        "latitude": estate["latitude"],
        "longitude": estate["longitude"],
        "address": estate["address"],
        "rent": estate["rent"],
        "doorHeight": estate["door_height"],
        "doorWidth": estate["door_width"],
        "popularity": estate["popularity"],
        "description": estate["description"],
        "features": estate["features"]
    }, ensure_ascii=False)


def generate_estate_dummy_data(estate_id, wrap={}):
    latlng = fake.local_latlng(country_code='JP', coords_only=True)
    feature_length = random.randint(0, min(3, len(ESTATE_FEATURE_LIST)))
    image_hash = fake.word(ext_word_list=ESTATE_IMAGE_HASH_LIST)

    estate = {
        "id": estate_id,
        "thumbnail": f'/images/estate/{image_hash}.png',
        "name": fake.word(ext_word_list=BUILDING_NAME_LIST).format(name=fake.last_name()),
        "latitude": float(latlng[0]) + random.normalvariate(mu=0.0, sigma=0.3),
        "longitude": float(latlng[1]) + random.normalvariate(mu=0.0, sigma=0.3),
        "address": fake.address(),
        "rent": random.randint(30000, 200000),
        "door_height": random.randint(DOOR_MIN_CENTIMETER, DOOR_MAX_CENTIMETER),
        "door_width": random.randint(DOOR_MIN_CENTIMETER, DOOR_MAX_CENTIMETER),
        "popularity": random.randint(MIN_POPULARITY, MAX_POPULARITY),
        "description": random.choice(desc_lines).strip(),
        "features": ','.join(fake.words(nb=feature_length, ext_word_list=ESTATE_FEATURE_LIST, unique=True))
    }
    return dict(estate, **wrap)


if __name__ == '__main__':

    for i, random_hash in enumerate(ESTATE_IMAGE_HASH_LIST):
        image_data_list = [read_src_file_data(
            image) for image in glob.glob(os.path.join(ESTATE_IMAGE_ORIGIN_DIR, "*.png"))]
        with open(os.path.join(ESTATE_IMAGE_PUBLIC_DIR, f"{random_hash}.png"), mode='wb') as image_file:
            image_file.write(
                image_data_list[i % len(image_data_list)] + random_hash.encode('utf-8'))

    with open(DESCRIPTION_LINES_FILE, mode='r', encoding='utf-8') as description_lines:
        desc_lines = description_lines.readlines()

    estate_id = 1

    with open(OUTPUT_SQL_FILE, mode='w', encoding='utf-8') as sqlfile, open(OUTPUT_TXT_FILE, mode='w', encoding='utf-8') as txtfile:
        if RECORD_COUNT % BULK_INSERT_COUNT != 0:
            raise Exception("The results of RECORD_COUNT and BULK_INSERT_COUNT need to be a divisible number. RECORD_COUNT = {}, BULK_INSERT_COUNT = {}".format(
                RECORD_COUNT, BULK_INSERT_COUNT))

        for _ in range(RECORD_COUNT//BULK_INSERT_COUNT):
            bulk_list = [generate_estate_dummy_data(
                estate_id + i) for i in range(BULK_INSERT_COUNT)]
            estate_id += BULK_INSERT_COUNT
            sqlCommand = f"""INSERT INTO isuumo.estate (id, thumbnail, name, latitude, longitude, address, rent, door_height, door_width, popularity, description, features) VALUES {', '.join(map(lambda estate: f"('{estate['id']}', '{estate['thumbnail']}', '{estate['name']}', '{estate['latitude']}' , '{estate['longitude']}', '{estate['address']}', '{estate['rent']}', '{estate['door_height']}', '{estate['door_width']}', '{estate['popularity']}', '{estate['description']}', '{estate['features']}')", bulk_list))};"""
            sqlfile.write(sqlCommand)
            txtfile.write("\n".join([dump_estate_to_json_str(estate)
                                     for estate in bulk_list]) + "\n")

    with open(VERIFY_DRAFT_FILE, mode='w', encoding='utf-8') as verify_draft_file:
        verify_draft_estates = [generate_estate_dummy_data(
            estate_id + i) for i in range(VERIFY_DRAFT_COUNT)]
        estate_id += VERIFY_DRAFT_COUNT
        verify_draft_file.write(
            "\n".join([dump_estate_to_json_str(estate) for estate in verify_draft_estates]) + "\n")

    for i in range(DRAFT_FILE_COUNT):
        with open(OUTPUT_DRAFT_FILE.format(index=i), mode='w', encoding='utf-8') as draft_file:
            draft_estates = [generate_estate_dummy_data(
                estate_id + i) for i in range(DRAFT_COUNT_PER_FILE)]
            estate_id += DRAFT_COUNT_PER_FILE
            draft_file.write(
                "\n".join([dump_estate_to_json_str(estate) for estate in draft_estates]) + "\n")

    with open(OUTPUT_FIXTURE_FILE, mode='w', encoding='utf-8') as fixture_file:
        fixture_file.write(json.dumps({
            "doorWidth": {
                "prefix": "",
                "suffix": "cm",
                "ranges": generate_ranges_from_separator(DOOR_WIDTH_RANGE_SEPARATORS)
            },
            "doorHeight": {
                "prefix": "",
                "suffix": "cm",
                "ranges": generate_ranges_from_separator(DOOR_HEIGHT_RANGE_SEPARATORS)
            },
            "rent": {
                "prefix": "",
                "suffix": "円",
                "ranges": generate_ranges_from_separator(RENT_RANGE_SEPARATORS)
            },
            "feature": {
                "list": ESTATE_FEATURE_LIST
            }
        }, ensure_ascii=False, indent=2))
