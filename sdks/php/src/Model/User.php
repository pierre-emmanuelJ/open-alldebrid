<?php

// Generated from openapi.yaml by scripts/gen_models.py. DO NOT EDIT.

declare(strict_types=1);

namespace OpenAllDebrid\Model;

use OpenAllDebrid\Internal\Hydrate;

final class User
{
    public function __construct(
        public readonly string $username,
        public readonly string $email,
        public readonly bool $isPremium,
        /** Whether the user has an active recurring subscription. */
        public readonly bool $isSubscribed,
        /** Whether the account is in its free trial. */
        public readonly bool $isTrial,
        /** Unix time when premium ends, or 0 when not premium. */
        public readonly int $premiumUntil,
        /** Interface language, e.g. `en` or `fr`. */
        public readonly string $lang,
        /** Preferred alldebrid.* top-level domain, e.g. `com` or `fr`. */
        public readonly string $preferedDomain,
        public readonly int $fidelityPoints,
        /** @var array<string, int> */
        public readonly array $limitedHostersQuotas,
        /**
         * Codes of the current notifications; clear them with `/v4/user/notification/clear`.
         * @var list<string>
         */
        public readonly array $notifications,
        /** During the trial, the remaining traffic in MB. */
        public readonly ?int $remainingTrialQuota = null,
        /** Whether recent links history logging is enabled. */
        public readonly ?bool $history = null,
    ) {
    }

    /** @param array<array-key, mixed> $data A decoded JSON object. */
    public static function fromArray(array $data): self
    {
        return new self(
            username: Hydrate::string(Hydrate::required($data, 'username')),
            email: Hydrate::string(Hydrate::required($data, 'email')),
            isPremium: Hydrate::bool(Hydrate::required($data, 'isPremium')),
            isSubscribed: Hydrate::bool(Hydrate::required($data, 'isSubscribed')),
            isTrial: Hydrate::bool(Hydrate::required($data, 'isTrial')),
            premiumUntil: Hydrate::int(Hydrate::required($data, 'premiumUntil')),
            lang: Hydrate::string(Hydrate::required($data, 'lang')),
            preferedDomain: Hydrate::string(Hydrate::required($data, 'preferedDomain')),
            fidelityPoints: Hydrate::int(Hydrate::required($data, 'fidelityPoints')),
            limitedHostersQuotas: Hydrate::map(Hydrate::emptyListAsMap(Hydrate::required($data, 'limitedHostersQuotas')), static fn (mixed $v): int => Hydrate::int($v)),
            notifications: Hydrate::list(Hydrate::required($data, 'notifications'), static fn (mixed $v): string => Hydrate::string($v)),
            remainingTrialQuota: array_key_exists('remainingTrialQuota', $data) && $data['remainingTrialQuota'] !== null ? Hydrate::int($data['remainingTrialQuota']) : null,
            history: array_key_exists('history', $data) && $data['history'] !== null ? Hydrate::bool($data['history']) : null,
        );
    }
}
