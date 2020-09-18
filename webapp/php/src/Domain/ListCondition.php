<?php

namespace App\Domain;

class ListCondition
{
    /** @var string[] */
    public array $list;

    public function __construct(array $list)
    {
        $this->list = $list;
    }

    public static function unmarshal(array $json): ListCondition
    {
        return new ListCondition($json['list']);
    }
}