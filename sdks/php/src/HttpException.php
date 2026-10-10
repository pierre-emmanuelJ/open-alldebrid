<?php

declare(strict_types=1);

namespace OpenAllDebrid;

use RuntimeException;

/** An HTTP status other than 200, such as 429 when rate limited. */
final class HttpException extends RuntimeException
{
    public function __construct(
        public readonly int $status,
        /** The API error in the body, when there is one. */
        public readonly ?AllDebridException $apiError = null,
    ) {
        parent::__construct($apiError !== null ? "HTTP {$status}: {$apiError->getMessage()}" : "HTTP {$status}");
    }

    /** Whether the API asked to slow down (HTTP 429 or 503). */
    public function isRateLimited(): bool
    {
        return $this->status === 429 || $this->status === 503;
    }
}
