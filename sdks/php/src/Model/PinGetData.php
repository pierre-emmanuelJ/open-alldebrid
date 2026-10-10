<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class PinGetData
{
    public function __construct(
        /** The code to show to the user. */
        public readonly string $pin,
        /** The token to send to `/v4/pin/check`. */
        public readonly string $check,
        /** Seconds before the PIN expires. */
        public readonly int $expiresIn,
        /** The page where the user enters the PIN, with the PIN filled in. */
        public readonly string $userUrl,
        /** The page where the user enters the PIN. */
        public readonly string $baseUrl,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            pin: Hydrate::string(Hydrate::required($data, 'pin')),
            check: Hydrate::string(Hydrate::required($data, 'check')),
            expiresIn: Hydrate::int(Hydrate::required($data, 'expires_in')),
            userUrl: Hydrate::string(Hydrate::required($data, 'user_url')),
            baseUrl: Hydrate::string(Hydrate::required($data, 'base_url')),
        );
    }
}
