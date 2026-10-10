<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class SavedLinksData
{
    public function __construct(
        /** @var list<SavedLink> */
        public readonly array $links,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            links: Hydrate::list(Hydrate::required($data, 'links'), static fn (mixed $v): SavedLink => SavedLink::fromArray(Hydrate::object($v))),
        );
    }
}
