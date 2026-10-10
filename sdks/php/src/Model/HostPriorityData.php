<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class HostPriorityData
{
    public function __construct(
        /**
         * Host names mapped to their rank, 1 being the most open.
         * @var array<string, int>
         */
        public readonly array $hosts,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            hosts: Hydrate::map(Hydrate::required($data, 'hosts'), static fn (mixed $v): int => Hydrate::int($v)),
        );
    }
}
