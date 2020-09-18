<?php

namespace App\Domain;

class Range
{
    public ?int $id;
    public ?int $min;
    public ?int $max;

    public function __construct(
        int $id = null,
        int $min = null,
        int $max = null
    ) {
        $this->id = $id;
        $this->min = $min;
        $this->max = $max;
    }

    public static function unmarshal(array $json) {
        return new Range(
            $json['id'] ?? null,
            $json['min'] ?? null,
            $json['max'] ?? null,
        );
    }
}