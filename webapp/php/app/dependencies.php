<?php
declare(strict_types=1);

use DI\ContainerBuilder;
use Monolog\Handler\StreamHandler;
use Monolog\Logger;
use Monolog\Processor\UidProcessor;
use Psr\Container\ContainerInterface;
use Psr\Log\LoggerInterface;

return function (ContainerBuilder $containerontainerBuilder) {
    $containerontainerBuilder->addDefinitions([
        'logger' => function(ContainerInterface $container): LoggerInterface {
            return $container->get(LoggerInterface::class);
        }
    ]);

    $containerontainerBuilder->addDefinitions([
        LoggerInterface::class => function(ContainerInterface $container): LoggerInterface {
            $settings = $container->get('settings');

            $loggerSettings = $settings['logger'];
            $logger = new Logger($loggerSettings['name']);

            $processor = new UidProcessor();
            $logger->pushProcessor($processor);

            $handler = new StreamHandler($loggerSettings['path'], $loggerSettings['level']);
            $logger->pushHandler($handler);

            return $logger;
        }
    ]);

    $containerontainerBuilder->addDefinitions([
        PDO::class => function(ContainerInterface $container): PDO {
            $settings = $container->get('settings')['database'];

            $dsn = vsprintf('mysql:host=%s;dbname=%s;port=%d', [
                $settings['host'],
                $settings['dbname'],
                $settings['port']
            ]);

            $pdo = new PDO($dsn, $settings['user'], $settings['pass']);
            $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
            $pdo->setAttribute(PDO::ATTR_DEFAULT_FETCH_MODE, PDO::FETCH_ASSOC);

            return $pdo;
        }
    ]);
};
