<?php
declare(strict_types=1);

use App\Domain\BoundingBox;
use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Fig\Http\Message\StatusCodeInterface;
use League\Csv\Reader;
use Slim\App;
use App\Domain\Chair;
use App\Domain\ChairSearchCondition;
use App\Domain\Coordinate;
use App\Domain\Estate;
use App\Domain\EstateSearchCondition;
use App\Domain\Range;
use App\Domain\RangeCondition;

const EXEC_SUCCESS = 0;
const NUM_LIMIT = 20;
const NUM_NAZOTTE_LIMIT = 50;

function getRange(RangeCondition $condition, string $rangeId): ?Range
{
    $id = intval($rangeId, 10);

    foreach ($condition->ranges as $range) {
        if ($range->id == $id) {
            return $range;
        }
    }

    return null;
}

return function (App $app) {
    $app->options('/{routes:.*}', function (Request $request, Response $response) {
        // CORS Pre-Flight OPTIONS Request Handler
        return $response;
    });

    $app->post('/initialize', function(Request $request, Response $response): Response {
        $config = $this->get('settings')['database'];

        $paths = [
            '../../mysql/db/0_Schema.sql',
            '../../mysql/db/1_DummyEstateData.sql',
            '../../mysql/db/2_DummyChairData.sql',
        ];

        foreach ($paths as $path) {
            $sqlFile = realpath($path);
            $cmdStr = vsprintf('mysql -h %s -u %s -p%s -P %s %s < %s', [
                $config['host'],
                $config['user'],
                $config['pass'],
                $config['port'],
                $config['dbname'],
                $sqlFile,
            ]);

            system("bash -c \"$cmdStr\"", $result);
            if ($result !== EXEC_SUCCESS) {
                $this->get('logger')->error('Initialize script error');
                return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
            }
        }

        $response->getBody()->write(json_encode([
            'language' => 'php',
        ]));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/chair/search', function(Request $request, Response $response) {
        $conditions = [];
        $params = [];

        /** @var ChairSearchCondition */
        $chairSearchCondition = $this->get(ChairSearchCondition::class);

        $priceRangeId = $request->getQueryParams()['priceRangeId'] ?? null;
        if (is_numeric($priceRangeId)) {
            if (!$chairPrice = getRange($chairSearchCondition->price, $priceRangeId)) {
                $this->get('logger')->info(sprintf('priceRangeId invalid, %s', $priceRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($chairPrice->min != -1) {
                $conditions[] = 'price >= :minPrice';
                $params[':minPrice'] = [$chairPrice->min, PDO::PARAM_INT];
            }
            if ($chairPrice->max != -1) {
                $conditions[] = 'price < :maxPrice';
                $params[':maxPrice'] = [$chairPrice->max, PDO::PARAM_INT];
            }
        }
        $heightRangeId = $request->getQueryParams()['heightRangeId'] ?? null;
        if (is_numeric($heightRangeId)) {
            if (!$chairHeight = getRange($chairSearchCondition->height, $heightRangeId)) {
                $this->get('logger')->info(sprintf('heightRangeId invalid, %s', $heightRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($chairHeight->min != -1) {
                $conditions[] = 'height >= :minHeight';
                $params[':minHeight'] = [$chairHeight->min, PDO::PARAM_INT];
            }
            if ($chairHeight->max != -1) {
                $conditions[] = 'height < :maxHeight';
                $params[':maxHeight'] = [$chairHeight->max, PDO::PARAM_INT];
            }
        }
        $widthRangeId = $request->getQueryParams()['widthRangeId'] ?? null;
        if (is_numeric($widthRangeId)) {
            if (!$chairWidth = getRange($chairSearchCondition->width, $widthRangeId)) {
                $this->get('logger')->info(sprintf('widthRangeId invalid, %s', $widthRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($chairWidth->min != -1) {
                $conditions[] = 'width >= :minWidth';
                $params[':minWidth'] = [$chairWidth->min, PDO::PARAM_INT];
            }
            if ($chairWidth->max != -1) {
                $conditions[] = 'width < :maxWidth';
                $params[':maxWidth'] = [$chairWidth->max, PDO::PARAM_INT];
            }
        }
        $depthRangeId = $request->getQueryParams()['depthRangeId'] ?? null;
        if (is_numeric($depthRangeId)) {
            if (!$chairDepth = getRange($chairSearchCondition->depth, $depthRangeId)) {
                $this->get('logger')->info(sprintf('depthRangeId invalid, %s', $depthRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($chairDepth->min != -1) {
                $conditions[] = 'depth >= :minDepth';
                $params[':minDepth'] = [$chairDepth->min, PDO::PARAM_INT];
            }
            if ($chairDepth->max != -1) {
                $conditions[] = 'depth < :maxDepth';
                $params[':maxDepth'] = [$chairDepth->max, PDO::PARAM_INT];
            }
        }
        if ($kind = $request->getQueryParams()['kind'] ?? null) {
            $conditions[] = 'kind = :kind';
            $params[':kind'] = [$kind, PDO::PARAM_STR];
        }
        if ($color = $request->getQueryParams()['color'] ?? null) {
            $conditions[] = 'color = :color';
            $params[':color'] = [$color, PDO::PARAM_STR];
        }
        if ($features = $request->getQueryParams()['features'] ?? null) {
            foreach (explode(',', $features) as $key => $feature) {
                $name = sprintf(':feature_%s', $key);
                $conditions[] = sprintf("features LIKE CONCAT('%%', %s, '%%')", $name);
                $params[$name] = [$feature, PDO::PARAM_STR];
            }
        }

        if (count($conditions) === 0) {
            $this->get('logger')->info('Search condition not found');
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $conditions[] = 'stock > 0';

        if (is_null($page = $request->getQueryParams()['page'] ?? null)) {
            $this->get('logger')->info(sprintf('Invalid format page parameter: %s', $page));
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }
        if (is_null($perPage = $request->getQueryParams()['perPage'] ?? null)) {
            $this->get('logger')->info(sprintf('Invalid format perPage parameter: %s', $perPage));
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $searchQuery = 'SELECT * FROM chair WHERE ';
        $countQuery = 'SELECT COUNT(*) FROM chair WHERE ';
        $searchCondition = implode(' AND ', $conditions);
        $limitOffset = ' ORDER BY popularity DESC, id ASC LIMIT :limit OFFSET :offset';

        $stmt = $this->get(PDO::class)->prepare($countQuery . $searchCondition);
        foreach ($params as $key => $bind) {
            list($value, $type) = $bind;
            $stmt->bindValue($key, $value, $type);
        }
        $stmt->execute();
        $count = (int)$stmt->fetchColumn();

        $params[':limit'] = [(int)$perPage, PDO::PARAM_INT];
        $params[':offset'] = [(int)$page*$perPage, PDO::PARAM_INT];

        $stmt = $this->get(PDO::class)->prepare($searchQuery . $searchCondition . $limitOffset);
        foreach ($params as $key => $bind) {
            list($value, $type) = $bind;
            $stmt->bindValue($key, $value, $type);
        }
        $stmt->execute();
        $chairs = $stmt->fetchAll(PDO::FETCH_CLASS, Chair::class);

        $response->getBody()->write(json_encode([
            'count' => $count,
            'chairs' => array_map(
                function(Chair $chair) {
                    return $chair->toArray();
                },
                $chairs
            ),
        ]));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/chair/low_priced', function(Request $request, Response $response) {
        $query = 'SELECT * FROM chair WHERE stock > 0 ORDER BY price ASC, id ASC LIMIT :limit';
        $stmt = $this->get(PDO::class)->prepare($query);
        $stmt->bindValue(':limit', NUM_LIMIT, PDO::PARAM_INT);
        $stmt->execute();
        $chairs = $stmt->fetchAll(PDO::FETCH_CLASS, Chair::class);

        if (count($chairs) === 0) {
            $this->get('logger')->error('getLowPricedChair not found');
            $response->getBody()->write(json_encode([
                'chairs' => []
            ]));
            return $response->withHeader('Content-Type', 'application/json');
        }

        $response->getBody()->write(json_encode([
            'chairs' => array_map(
                function(Chair $chair) {
                    return $chair->toArray();
                },
                $chairs
            )
        ]));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/chair/search/condition', function(Request $request, Response $response) {
        $chairSearchCondition = $this->get(ChairSearchCondition::class);
        $response->getBody()->write(json_encode($chairSearchCondition));
        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->post('/api/chair/buy/{id}', function(Request $request, Response $response, array $args) {
        $id = $args['id'] ?? null;
        if (empty($id) || !is_numeric($id)) {
            $this->get('logger')->info('post request document failed');
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $pdo = $this->get(PDO::class);

        try {
            $pdo->beginTransaction();
            $stmt = $pdo->prepare('SELECT * FROM chair WHERE id = :id AND stock > 0 FOR UPDATE');
            $stmt->bindValue(':id', $id, PDO::PARAM_INT);
            $stmt->execute();
            $chair = $stmt->fetchObject(Chair::class);

            if (!$chair) {
                $pdo->rollBack();
                $this->get('logger')->info(sprintf('buyChair chair id "%s" not found', $id));
                return $response->withStatus(StatusCodeInterface::STATUS_NOT_FOUND);
            }

            $stmt = $pdo->prepare('UPDATE chair SET stock = stock - 1 WHERE id = :id');
            $stmt->bindValue(':id', $id, PDO::PARAM_INT);
            if (!$stmt->execute()) {
                $pdo->rollBack();
                $this->get('logger')->error('chair stock update failed');
                return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
            }

            $pdo->commit();
        } catch (PDOException $e) {
            $pdo->rollBack();
            $this->get('logger')->error(sprintf('DB Execution Error: on getting a chair by id : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/chair/{id}', function(Request $request, Response $response, array $args) {
        $id = $args['id'] ?? null;
        if (empty($id) || !is_numeric($id)) {
            $this->get('logger')->error(sprintf('Request parameter \"id\" parse error : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $query = 'SELECT * FROM chair WHERE id = :id';
        $stmt = $this->get(PDO::class)->prepare($query);
        $stmt->execute([':id' => $id]);
        $chair = $stmt->fetchObject(Chair::class);

        if (!$chair) {
            $this->get('logger')->error(sprintf('requested id\'s chair not found : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_NOT_FOUND);
        } elseif (!$chair instanceof Chair) {
            $this->get('logger')->error(sprintf('Failed to get the chair from id : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        } elseif ($chair->getStock() <= 0) {
            $this->get('logger')->error(sprintf('requested id\'s chair is sold out : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_NOT_FOUND);
        }

        $response->getBody()->write(json_encode($chair->toArray()));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->post('/api/chair', function (Request $request, Response $response) {
        if (!$file = $request->getUploadedFiles()['chairs'] ?? null) {
            $this->get('logger')->error('failed to get form file');
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        } elseif (!$file instanceof Slim\Psr7\UploadedFile || $file->getError() !== UPLOAD_ERR_OK) {
            $this->get('logger')->error('failed to get form file');
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        if (!$records = Reader::createFromPath($file->getFilePath())) {
            $this->get('logger')->error('failed to read csv');
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        $pdo = $this->get(PDO::class);

        try {
            $pdo->beginTransaction();

            foreach ($records as $record) {
                $query = 'INSERT INTO chair VALUES(:id, :name, :description, :thumbnail, :price, :height, :width, :depth, :color, :features, :kind, :popularity, :stock)';
                $stmt = $pdo->prepare($query);
                $stmt->execute([
                    ':id' => (int)trim($record[0] ?? null),
                    ':name' => (string)trim($record[1] ?? null),
                    ':description' => (string)trim($record[2] ?? null),
                    ':thumbnail' => (string)trim($record[3] ?? null),
                    ':price' => (int)trim($record[4] ?? null),
                    ':height' => (int)trim($record[5] ?? null),
                    ':width' => (int)trim($record[6]) ?? null,
                    ':depth' => (int)trim($record[7] ?? null),
                    ':color' => (string)trim($record[8] ?? null),
                    ':features' => (string)trim($record[9] ?? null),
                    ':kind' => (string)trim($record[10] ?? null),
                    ':popularity' => (int)trim($record[11] ?? null),
                    ':stock' => (int)trim($record[12] ?? null),
                ]);
            }

            $pdo->commit();
        } catch (PDOException $e) {
            $pdo->rollBack();
            $this->get('logger')->error(sprintf('failed to insert chair: %s', $e->getMessage()));
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        return $response->withStatus(StatusCodeInterface::STATUS_CREATED);
    });

    // Estate
    $app->post('/api/estate', function(Request $request, Response $response) {
        if (!$file = $request->getUploadedFiles()['estates'] ?? null) {
            $this->get('logger')->error('failed to get form file');
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        } elseif (!$file instanceof Slim\Psr7\UploadedFile || $file->getError() !== UPLOAD_ERR_OK) {
            $this->get('logger')->error('failed to get form file');
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        if (!$records = Reader::createFromPath($file->getFilePath())) {
            $this->get('logger')->error('failed to read csv');
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        $pdo = $this->get(PDO::class);

        try {
            $pdo->beginTransaction();

            foreach ($records as $record) {
                $query = 'INSERT INTO estate VALUES(:id, :name, :description, :thumbnail, :address, :latitude, :longitude, :rent, :door_height, :door_width, :features, :popularity)';
                $stmt = $pdo->prepare($query);
                $stmt->execute([
                    'id' => (int)trim($record[0] ?? null),
                    'name' => trim($record[1] ?? null),
                    'description' => trim($record[2] ?? null),
                    'thumbnail' => trim($record[3] ?? null),
                    'address' => trim($record[4] ?? null),
                    'latitude' => (float)trim($record[5] ?? null),
                    'longitude' => (float)trim($record[6] ?? null),
                    'rent' => (int)trim($record[7] ?? null),
                    'door_height' => (int)trim($record[8] ?? null),
                    'door_width' => (int)trim($record[9] ?? null),
                    'features' => trim($record[10] ?? null),
                    'popularity' => (int)trim($record[11] ?? null),
                ]);
            }

            $pdo->commit();
        } catch (PDOException $e) {
            $pdo->rollBack();
            $this->get('logger')->error(sprintf('failed to insert estate: %s', $e->getMessage()));
            return $response->withStatus(StatusCodeInterface::STATUS_INTERNAL_SERVER_ERROR);
        }

        return $response->withStatus(StatusCodeInterface::STATUS_CREATED);
    });

    $app->get('/api/estate/search', function(Request $request, Response $response) {
        $conditions = [];
        $params = [];

        /** @var EstateSearchCondition */
        $estateSearchCondition = $this->get(EstateSearchCondition::class);

        $doorHeightRangeId = $request->getQueryParams()['doorHeightRangeId'] ?? null;
        if (is_numeric($doorHeightRangeId)) {
            if (!$doorHeight = getRange($estateSearchCondition->doorHeight, $doorHeightRangeId)) {
                $this->get('logger')->info(sprintf('doorHeightRangeId invalid, %s', $doorHeightRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($doorHeight->min != -1) {
                $conditions[] = 'door_height >= :minDoorHeight';
                $params[':minDoorHeight'] = [$doorHeight->min, PDO::PARAM_INT];
            }
            if ($doorHeight->max != -1) {
                $conditions[] = 'door_height < :maxDoorHeight';
                $params[':maxDoorHeight'] = [$doorHeight->max, PDO::PARAM_INT];
            }
        }
        $doorWidthRangeId = $request->getQueryParams()['doorWidthRangeId'] ?? null;
        if (is_numeric($doorWidthRangeId)) {
            if (!$doorWidth = getRange($estateSearchCondition->doorWidth, $doorWidthRangeId)) {
                $this->get('logger')->info(sprintf('doorWidthRangeId invalid, %s', $doorWidthRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($doorWidth->min != -1) {
                $conditions[] = 'door_width >= :minDoorWidth';
                $params[':minDoorWidth'] = [$doorWidth->min, PDO::PARAM_INT];
            }
            if ($doorWidth->max != -1) {
                $conditions[] = 'door_width < :maxDoorWidth';
                $params[':maxDoorWidth'] = [$doorWidth->max, PDO::PARAM_INT];
            }
        }
        $rentRangeId = $request->getQueryParams()['rentRangeId'] ?? null;
        if (is_numeric($rentRangeId)) {
            if (!$estateRent = getRange($estateSearchCondition->rent, $rentRangeId)) {
                $this->get('logger')->info(sprintf('rentRangeId invalid, %s', $rentRangeId));
                return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
            }
            if ($estateRent->min != -1) {
                $conditions[] = 'rent >= :minEstateRent';
                $params[':minEstateRent'] = [$estateRent->min, PDO::PARAM_INT];
            }
            if ($estateRent->max != -1) {
                $conditions[] = 'rent < :maxEstateRent';
                $params[':maxEstateRent'] = [$estateRent->max, PDO::PARAM_INT];
            }
        }

        if ($features = $request->getQueryParams()['features'] ?? null) {
            foreach (explode(',', $features) as $key => $feature) {
                $name = sprintf(':feature_%s', $key);
                $conditions[] = sprintf("features LIKE CONCAT('%%', %s, '%%')", $name);
                $params[$name] = [$feature, PDO::PARAM_STR];
            }
        }

        if (count($conditions) === 0) {
            $this->get('logger')->info('Search condition not found');
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        if (is_null($page = $request->getQueryParams()['page'] ?? null)) {
            $this->get('logger')->info(sprintf('Invalid format page parameter: %s', $page));
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }
        if (is_null($perPage = $request->getQueryParams()['perPage'] ?? null)) {
            $this->get('logger')->info(sprintf('Invalid format perPage parameter: %s', $perPage));
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $searchQuery = 'SELECT * FROM estate WHERE ';
        $countQuery = 'SELECT COUNT(*) FROM estate WHERE ';
        $searchCondition = implode(' AND ', $conditions);
        $limitOffset = ' ORDER BY popularity DESC, id ASC LIMIT :limit OFFSET :offset';

        $stmt = $this->get(PDO::class)->prepare($countQuery . $searchCondition);
        foreach ($params as $key => $bind) {
            list($value, $type) = $bind;
            $stmt->bindValue($key, $value, $type);
        }
        $stmt->execute();
        $count = (int)$stmt->fetchColumn();

        $params[':limit'] = [(int)$perPage, PDO::PARAM_INT];
        $params[':offset'] = [(int)$page*$perPage, PDO::PARAM_INT];

        $stmt = $this->get(PDO::class)->prepare($searchQuery . $searchCondition . $limitOffset);
        foreach ($params as $key => $bind) {
            list($value, $type) = $bind;
            $stmt->bindValue($key, $value, $type);
        }
        $stmt->execute();
        $estates = $stmt->fetchAll(PDO::FETCH_CLASS, Estate::class);

        $response->getBody()->write(json_encode([
            'count' => $count,
            'estates' => array_map(
                function(Estate $estate) {
                    return $estate->toArray();
                },
                $estates
            )
        ]));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/estate/low_priced', function(Request $request, Response $response) {
        $query = 'SELECT * FROM estate ORDER BY rent ASC, id ASC LIMIT :limit';
        $stmt = $this->get(PDO::class)->prepare($query);
        $stmt->bindValue(':limit', NUM_LIMIT, PDO::PARAM_INT);
        $stmt->execute();
        $estates = $stmt->fetchAll(PDO::FETCH_CLASS, Estate::class);

        if (count($estates) === 0) {
            $this->get('logger')->error('getLowPricedEstate not found');
            $response->getBody()->write(json_encode([
                'chairs' => []
            ]));
            return $response->withHeader('Content-Type', 'application/json');
        }

        $response->getBody()->write(json_encode([
            'estates' => array_map(
                function(Estate $estate) {
                    return $estate->toArray();
                },
                $estates
            )
        ]));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->post('/api/estate/req_doc/{id}', function(Request $request, Response $response, array $args) {
        $id = $args['id'] ?? null;
        if (empty($id) || !is_numeric($id)) {
            $this->get('logger')->error(sprintf('Request parameter "id" parse error : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $query = 'SELECT * FROM estate WHERE id = :id';
        $stmt = $this->get(PDO::class)->prepare($query);
        $stmt->execute([':id' => $id]);
        $estate = $stmt->fetchObject(Estate::class);

        if (!$estate) {
            $this->get('logger')->error(sprintf('requested id\'s estate not found : %s', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_NOT_FOUND);
        }

        return $response->withStatus(StatusCodeInterface::STATUS_OK);
    });

    $app->post('/api/estate/nazotte', function(request $request, Response $response) {
        $json = json_decode($request->getBody()->getContents(), true);
        $coordinates = array_map(
            Coordinate::class . '::createFromJson',
            $json['coordinates']
        );
        if (count($coordinates) === 0) {
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $boundingBox = BoundingBox::createFromCordinates($coordinates);

        $query = 'SELECT * FROM estate WHERE latitude <= ? AND latitude >= ? AND longitude <= ? AND longitude >= ? ORDER BY popularity DESC, id ASC';
        $stmt = $this->get(PDO::class)->prepare($query);
        $stmt->execute([
            $boundingBox->bottomRightCorner->latitude,
            $boundingBox->topLeftCorner->latitude,
            $boundingBox->bottomRightCorner->longitude,
            $boundingBox->topLeftCorner->longitude,
        ]);
        $estatesInBoundingBox = $stmt->fetchAll(PDO::FETCH_CLASS, Estate::class);

        $estatesInPolygon = Array();
        foreach ($estatesInBoundingBox as $estate) {
            $point = sprintf("'POINT(%f %f)'", $estate->latitude, $estate->longitude);
            $query = sprintf('SELECT * FROM estate WHERE id = ? AND ST_Contains(ST_PolygonFromText(%s), ST_GeomFromText(%s))', Coordinate::toText($coordinates), $point);
            $stmt = $this->get(PDO::class)->prepare($query);
            $stmt->execute([$estate->id]);
            $e = $stmt->fetchObject(Estate::class);
            if ($e) {
                array_push($estatesInPolygon, $e);
            }
        }

        if (count($estatesInPolygon) > NUM_NAZOTTE_LIMIT) {
            $estatesInPolygon = array_slice($estatesInPolygon, 0, NUM_NAZOTTE_LIMIT);
        }

        $response->getBody()->write(json_encode([
            'count' => count($estatesInPolygon),
            'estates' => array_map(
                function(Estate $estate) {
                    return $estate->toArray();
                },
                $estatesInPolygon
            ),
        ]));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/estate/search/condition', function(Request $request, Response $response) {
        $estateSearchCondition = $this->get(EstateSearchCondition::class);
        $response->getBody()->write(json_encode($estateSearchCondition));
        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/recommended_estate/{id}', function(Request $request, Response $response, array $args) {
        $id = $args['id'] ?? null;
        if (empty($id) || !is_numeric($id)) {
            $this->get('logger')->info('post request document failed');
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $query = 'SELECT * FROM chair WHERE id = ?';
        $stmt = $this->get(PDO::class)->prepare($query);
        $stmt->execute([$id]);
        $chair = $stmt->fetchObject(Chair::class);

        if (!$chair) {
            $this->get('logger')->info(sprintf('Requested chair id "%s" not found', $id));
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $query = 'SELECT * FROM estate WHERE (door_width >= :w AND door_height >= :h) OR (door_width >= :w AND door_height >= :d) OR (door_width >= :h AND door_height >= :w) OR (door_width >= :h AND door_height >= :d) OR (door_width >= :d AND door_height >= :w) OR (door_width >= :d AND door_height >= :h) ORDER BY popularity DESC, id ASC LIMIT :limit';
        $stmt = $this->get(PDO::class)->prepare($query);
        $stmt->bindValue(':w', $chair->getWidth(), PDO::PARAM_INT);
        $stmt->bindValue(':h', $chair->getHeight(), PDO::PARAM_INT);
        $stmt->bindValue(':d', $chair->getDepth(), PDO::PARAM_INT);
        $stmt->bindValue(':limit', NUM_LIMIT, PDO::PARAM_INT);
        $stmt->execute();
        $estates = $stmt->fetchAll(PDO::FETCH_CLASS, Estate::class);

        $response->getBody()->write(json_encode([
            'estates' => array_map(
                function(Estate $estate) {
                    return $estate->toArray();
                },
                $estates
            ),
        ]));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/api/estate/{id}', function(Request $request, Response $response, array $args) {
        $id = $args['id'] ?? null;
        if (empty($id) || !is_numeric($id)) {
            $this->get('logger')->info('Request parameter "id" parse error');
            return $response->withStatus(StatusCodeInterface::STATUS_BAD_REQUEST);
        }

        $stmt = $this->get(PDO::class)->prepare('SELECT * FROM estate WHERE id = :id');
        $stmt->bindValue(':id', $id, PDO::PARAM_INT);
        $stmt->execute();

        $estate = $stmt->fetchObject(Estate::class);

        $response->getBody()->write(json_encode($estate->toArray()));

        return $response->withHeader('Content-Type', 'application/json');
    });

    $app->get('/', function (Request $request, Response $response) {
        $response->getBody()->write('Hello world!');
        return $response;
    });
};
