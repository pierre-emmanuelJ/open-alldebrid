<?php

declare(strict_types=1);

namespace OpenAllDebrid;

use RuntimeException;

/** An error answered by the API ("status": "error"). Compare errorCode with ErrorCode constants. */
final class AllDebridException extends RuntimeException
{
    public function __construct(
        /** The error code, e.g. MAGNET_INVALID_ID. */
        public readonly string $errorCode,
        /** The message of the API. */
        public readonly string $apiMessage,
        /** With AUTH_BLOCKED, the token to poll with AllDebrid::verificationStatus(). */
        public readonly ?string $token = null,
    ) {
        parent::__construct("{$errorCode}: {$apiMessage}");
    }
}
