<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class VerificationData
{
    public function __construct(
        public readonly string $verif,
        /** When `waiting`, whether `/v4/user/verif/resend` may be called. */
        public readonly ?bool $resendable = null,
        /** When `allowed`, the API key. */
        public readonly ?string $apikey = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            verif: Hydrate::string(Hydrate::required($data, 'verif')),
            resendable: array_key_exists('resendable', $data) && $data['resendable'] !== null ? Hydrate::bool($data['resendable']) : null,
            apikey: array_key_exists('apikey', $data) && $data['apikey'] !== null ? Hydrate::string($data['apikey']) : null,
        );
    }
}
