<?php
declare(strict_types=1);

use App\Domain\ChairSearchCondition;
use App\Domain\EstateSearchCondition;
use DI\ContainerBuilder;
use Psr\Container\ContainerInterface;

return function (ContainerBuilder $containerBuilder) {
    $containerBuilder->addDefinitions([
        ChairSearchCondition::class => function (ContainerInterface $c) {
            if (!$jsonText = file_get_contents('../../fixture/chair_condition.json')) {
                throw new RuntimeException(sprintf('Failed to get load file: %s', '/fixture/chair_condition.json'));
            }
            if (!$json = json_decode($jsonText, true)) {
                throw new RuntimeException(sprintf('Failed to parse json: %s', '../fixture/chair_condition.json'));
            }
            return ChairSearchCondition::unmarshal($json);
        }
    ]);

    $containerBuilder->addDefinitions([
        EstateSearchCondition::class => function (ContainerInterface $c) {
            if (!$jsonText = file_get_contents('../../fixture/estate_condition.json')) {
                throw new RuntimeException(sprintf('Failed to get load file: %s', '/fixture/estate_condition.json'));
            }
            if (!$json = json_decode($jsonText, true)) {
                throw new RuntimeException(sprintf('Failed to parse json: %s', '../fixture/estate_condition.json'));
            }
            return EstateSearchCondition::unmarshal($json);
        }
    ]);
};
