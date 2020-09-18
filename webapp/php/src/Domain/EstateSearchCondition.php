<?php

namespace App\Domain;

class EstateSearchCondition
{
    public ?RangeCondition $doorWidth;
    public ?RangeCondition $doorHeight;
    public ?RangeCondition $rent;
    public ?ListCondition $feature;

    public function __construct(
        RangeCondition $doorHeight = null,
        RangeCondition $doorWidth = null,
        RangeCondition $rent = null,
        ListCondition $feature = null
    ) {
        $this->doorHeight = $doorHeight;
        $this->doorWidth = $doorWidth;
        $this->rent = $rent;
        $this->feature = $feature;
    }

    public static function unmarshal(array $json): EstateSearchCondition
    {
        return new EstateSearchCondition(
            isset($json['doorHeight']) ? RangeCondition::unmarshal($json['doorHeight']) : null,
            isset($json['doorWidth']) ? RangeCondition::unmarshal($json['doorWidth']) : null,
            isset($json['rent']) ? RangeCondition::unmarshal($json['rent']) : null,
            isset($json['feature']) ? ListCondition::unmarshal($json['feature']) : null,
        );
    }
}