<?php
declare(strict_types=1);

namespace App\Application\Middleware;

use Psr\Http\Message\ResponseInterface as Response;
use Psr\Http\Message\ServerRequestInterface as Request;
use Psr\Http\Server\MiddlewareInterface as Middleware;
use Psr\Http\Server\RequestHandlerInterface as RequestHandler;
use Psr\Log\LoggerInterface;

class LoggerMiddleware implements Middleware
{
    private LoggerInterface $logger;

    public function __construct(LoggerInterface $logger)
    {
        $this->logger = $logger;
    }

    /**
     * {@inheritdoc}
     */
    public function process(Request $request, RequestHandler $handler): Response
    {
        $start = microtime(true);

        $response = $handler->handle($request);

        $this->logger->warning(json_encode([
            'time' => microtime(true) - $start,
            'remote_ip' => $request->getServerParams()['REMOTE_ADDR'] ?? null,
            'host' => $request->getUri()->getHost(),
            'method' => $request->getMethod(),
            'uri' => $request->getUri()->getPath(),
            'user_agent' => $request->getHeader('User-Agent'),
            'status' => $response->getStatusCode(),
        ]));

        return $response;
    }
}
