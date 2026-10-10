<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class PinCheckData
{
    public function __construct(
        /** `true` once the user has entered the PIN. */
        public readonly bool $activated,
        /** Seconds before the PIN expires. */
        public readonly int $expiresIn,
        /** The new API key, once `activated` is `true`. */
        public readonly ?string $apikey = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            activated: Hydrate::bool(Hydrate::required($data, 'activated')),
            expiresIn: Hydrate::int(Hydrate::required($data, 'expires_in')),
            apikey: array_key_exists('apikey', $data) && $data['apikey'] !== null ? Hydrate::string($data['apikey']) : null,
        );
    }
}
