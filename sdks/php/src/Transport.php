<?php

declare(strict_types=1);

namespace OpenAllDebrid;

/**
 * Sends one HTTP request. The default, CurlTransport, uses ext-curl; pass
 * your own to use another HTTP client or to test.
 */
interface Transport
{
    /**
     * @param array<string, string> $headers
     * @return array{0: int, 1: string} The HTTP status and the body.
     */
    public function send(string $method, string $url, array $headers, ?string $body): array;
}
