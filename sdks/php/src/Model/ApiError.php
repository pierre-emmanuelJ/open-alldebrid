<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class ApiError
{
    public function __construct(
        /** A machine-readable error code; see `ErrorCode`. */
        public readonly string $code,
        /** A human-readable message, in English. */
        public readonly string $message,
        /** With `AUTH_BLOCKED`, the token to poll on `/v4/user/verif`. */
        public readonly ?string $token = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            code: Hydrate::string(Hydrate::required($data, 'code')),
            message: Hydrate::string(Hydrate::required($data, 'message')),
            token: array_key_exists('token', $data) && $data['token'] !== null ? Hydrate::string($data['token']) : null,
        );
    }
}
