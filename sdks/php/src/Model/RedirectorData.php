<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class RedirectorData
{
    public function __construct(
        /**
         * Encrypted links, usable with `/v4/link/infos` and `/v4/link/unlock`.
         * @var list<string>
         */
        public readonly array $links,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            links: Hydrate::list(Hydrate::required($data, 'links'), static fn (mixed $v): string => Hydrate::string($v)),
        );
    }
}
