export interface Estate {
  id: string
  name: string
  thumbnail: string
  address: string
  description: string
  doorHeight: number
  doorWidth: number
  features: string
  latitude: number
  longitude: number
  rent: number
}

export interface Chair {
  id: string
  name: string
  thumbnail: string
  description: string
  height: number
  width: number
  depth: number
  features: string
  price: number
  color: string
  kind: string
}

export interface Coordinate {
  latitude: number
  longitude: number
}

export interface Range {
  id: number
  min: number
  max: number
}

export interface RangeCondition {
  prefix: string
  suffix: string
  ranges: Range[]
}

export interface ListCondition {
  list: string[]
}

export interface EstateSearchCondition {
  doorWidth: RangeCondition
  doorHeight: RangeCondition
  rent: RangeCondition
  feature: ListCondition
}

export interface EstateSearchParams {
  doorWidthRangeId: string
  doorHeightRangeId: string
  rentRangeId: string
  features: string
  page: number
  perPage: number
}

export interface EstateSearchResponse {
  estates: Estate[]
  count: number
}

export interface ChairSearchCondition {
  price: RangeCondition
  height: RangeCondition
  width: RangeCondition
  depth: RangeCondition
  color: ListCondition
  feature: ListCondition
  kind: ListCondition
}

export interface ChairSearchParams {
  priceRangeId: string
  heightRangeId: string
  widthRangeId: string
  depthRangeId: string
  color: string
  kind: string
  features: string
  page: number
  perPage: number
}

export interface ChairSearchResponse {
  chairs: Chair[]
  count: number
}
