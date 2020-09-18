<?php

namespace App\Domain;

class RangeCondition
{
    public ?string $prefix;
    public ?string $suffix;
    /** @var Range[] */
    public array $ranges;

    public function __construct(
        string $prefix = null,
        string $suffix = null,
        array $ranges = []
    ) {
        $this->prefix = $prefix;
        $this->suffix = $suffix;
        $this->ranges = $ranges;
    }

    public static function unmarshal(array $json): RangeCondition
    {
        return new RangeCondition(
            $json['prefix'] ?? null,
            $json['suffix'] ?? null,
            array_map(Range::class . '::unmarshal', $json['ranges'] ?? [])
        );
    }
}