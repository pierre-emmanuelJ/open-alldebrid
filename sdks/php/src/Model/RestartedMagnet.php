<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class RestartedMagnet
{
    public function __construct(
        /** The magnet ID. */
        public readonly ?string $magnet = null,
        public readonly ?string $message = null,
        public readonly ?ApiError $error = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            magnet: array_key_exists('magnet', $data) && $data['magnet'] !== null ? Hydrate::string($data['magnet']) : null,
            message: array_key_exists('message', $data) && $data['message'] !== null ? Hydrate::string($data['message']) : null,
            error: array_key_exists('error', $data) && $data['error'] !== null ? ApiError::fromArray(Hydrate::object($data['error'])) : null,
        );
    }
}
