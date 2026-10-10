<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class UserHostsData
{
    public function __construct(
        /** @var array<string, Service> */
        public readonly array $hosts,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            hosts: Hydrate::map(Hydrate::required($data, 'hosts'), static fn (mixed $v): Service => Service::fromArray(Hydrate::object($v))),
        );
    }
}
