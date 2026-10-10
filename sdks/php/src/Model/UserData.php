<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class UserData
{
    public function __construct(
        public readonly User $user,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            user: User::fromArray(Hydrate::object(Hydrate::required($data, 'user'))),
        );
    }
}
