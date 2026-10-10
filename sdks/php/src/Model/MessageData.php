<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class MessageData
{
    public function __construct(
        /** A confirmation message. */
        public readonly string $message,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            message: Hydrate::string(Hydrate::required($data, 'message')),
        );
    }
}
