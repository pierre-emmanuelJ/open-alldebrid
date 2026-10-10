<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class VoucherBalanceData
{
    public function __construct(
        /** The reseller balance. */
        public readonly float $balance,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            balance: Hydrate::float(Hydrate::required($data, 'balance')),
        );
    }
}
