<?php

namespace App\Domain;

class ChairSearchCondition
{
    public ?RangeCondition $width;
    public ?RangeCondition $height;
    public ?RangeCondition $depth;
    public ?RangeCondition $price;
    public ?ListCondition $color;
    public ?ListCondition $feature;
    public ?ListCondition $kind;

    public function __construct(
        RangeCondition $width = null,
        RangeCondition $height = null,
        RangeCondition $depth = null,
        RangeCondition $price = null,
        ListCondition $color = null,
        ListCondition $feature = null,
        ListCondition $kind = null
    ) {
        $this->width = $width;
        $this->height = $height;
        $this->depth = $depth;
        $this->price = $price;
        $this->color = $color;
        $this->feature = $feature;
        $this->kind = $kind;
    }

    public static function unmarshal(array $json): ChairSearchCondition
    {
        return new ChairSearchCondition(
            isset($json['width']) ? RangeCondition::unmarshal($json['width']) : null,
            isset($json['height']) ? RangeCondition::unmarshal($json['height']) : null,
            isset($json['depth']) ? RangeCondition::unmarshal($json['depth']) : null,
            isset($json['price']) ? RangeCondition::unmarshal($json['price']) : null,
            isset($json['color']) ? ListCondition::unmarshal($json['color']) : null,
            isset($json['feature']) ? ListCondition::unmarshal($json['feature']) : null,
            isset($json['kind']) ? ListCondition::unmarshal($json['kind']) : null,
        );
    }
}