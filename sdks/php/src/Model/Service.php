<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

/** A supported file host, streaming site or redirector. */
final class Service
{
    public function __construct(
        public readonly string $name,
        /** `premium` (needs a premium account) or `free`. */
        public readonly string $type,
        /**
         * Every domain of the service.
         * @var list<string>
         */
        public readonly array $domains,
        /**
         * Regular expressions matching the supported links.
         * @var list<string>|null
         */
        public readonly ?array $regexps = null,
        /**
         * The `regexps` joined into one expression (sometimes a list).
         * @var list<string>|null
         */
        public readonly ?array $regexp = null,
        /**
         * Whether the service worked in the last 5 minutes. Only some hosts are
         * checked, about every 10 minutes.
         */
        public readonly ?bool $status = null,
        /** A display name, when it differs from `name`. */
        public readonly ?string $display = null,
        /**
         * Domains that redirect to this host.
         * @var list<string>|null
         */
        public readonly ?array $hardRedirect = null,
        /** Redirectors only. Extraction is slow or rate limited. */
        public readonly ?bool $heavy = null,
        /** `/v4.1/user/hosts` only. The remaining quota (see `quotaType`). */
        public readonly ?int $quota = null,
        /** `/v4.1/user/hosts` only. The daily quota. */
        public readonly ?int $quotaMax = null,
        /**
         * `/v4.1/user/hosts` only. `traffic` (quota in MB) or `nb_download`
         * (number of downloads).
         */
        public readonly ?string $quotaType = null,
        /** `/v4.1/user/hosts` only. Free simultaneous download slots. */
        public readonly ?int $limitSimuDl = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            name: Hydrate::string(Hydrate::required($data, 'name')),
            type: Hydrate::string(Hydrate::required($data, 'type')),
            domains: Hydrate::list(Hydrate::required($data, 'domains'), static fn (mixed $v): string => Hydrate::string($v)),
            regexps: array_key_exists('regexps', $data) && $data['regexps'] !== null ? Hydrate::list($data['regexps'], static fn (mixed $v): string => Hydrate::string($v)) : null,
            regexp: array_key_exists('regexp', $data) && $data['regexp'] !== null ? Hydrate::strings($data['regexp']) : null,
            status: array_key_exists('status', $data) && $data['status'] !== null ? Hydrate::bool($data['status']) : null,
            display: array_key_exists('display', $data) && $data['display'] !== null ? Hydrate::string($data['display']) : null,
            hardRedirect: array_key_exists('hardRedirect', $data) && $data['hardRedirect'] !== null ? Hydrate::list($data['hardRedirect'], static fn (mixed $v): string => Hydrate::string($v)) : null,
            heavy: array_key_exists('heavy', $data) && $data['heavy'] !== null ? Hydrate::bool($data['heavy']) : null,
            quota: array_key_exists('quota', $data) && $data['quota'] !== null ? Hydrate::int($data['quota']) : null,
            quotaMax: array_key_exists('quotaMax', $data) && $data['quotaMax'] !== null ? Hydrate::int($data['quotaMax']) : null,
            quotaType: array_key_exists('quotaType', $data) && $data['quotaType'] !== null ? Hydrate::string($data['quotaType']) : null,
            limitSimuDl: array_key_exists('limitSimuDl', $data) && $data['limitSimuDl'] !== null ? Hydrate::int($data['limitSimuDl']) : null,
        );
    }
}
