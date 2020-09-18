<?php
declare(strict_types=1);

use App\Application\Middleware\SessionMiddleware;
use App\Application\Middleware\LoggerMiddleware;
use Slim\App;

return function (App $app) {
    $container = $app->getContainer();
    $app->add(SessionMiddleware::class);
    $app->add(new LoggerMiddleware($container->get('logger')));
};
