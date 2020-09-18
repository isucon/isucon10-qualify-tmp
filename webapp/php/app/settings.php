<?php
declare(strict_types=1);

use DI\ContainerBuilder;
use Monolog\Logger;

return function (ContainerBuilder $containerBuilder) {
    // Global Settings Object
    $containerBuilder->addDefinitions([
        'settings' => [
            'displayErrorDetails' => true, // Should be set to false in production
            'logger' => [
                'name' => 'slim-app',
                'path' => 'php://stdout', // __DIR__ . '/var/log/app.log'
                'level' => Logger::DEBUG,
            ],
            'database' => [
                'host' => getenv('MYSQL_HOST') ?: '127.0.0.1',
                'port' => getenv('MYSQL_PORT') ?: '3306',
                'user' => getenv('MYSQL_USER') ?: 'isucon',
                'pass' => getenv('MYSQL_PASS') ?: 'isucon',
                'dbname' => getenv('MYSQL_DBNAME') ?: 'isuumo',
            ],
        ],
    ]);
};
