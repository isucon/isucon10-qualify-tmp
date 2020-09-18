import {
  FormControl,
  FormLabel,
  FormControlLabel,
  Radio,
  RadioGroup
} from '@material-ui/core'

import type { FC, ChangeEvent } from 'react'
import type { RangeCondition } from '@types'

interface Props {
  name: string
  value: string
  rangeCondition: RangeCondition
  onChange: (event: ChangeEvent<HTMLInputElement>, value: string) => void
}

export const RangeForm: FC<Props> = ({ name, value, rangeCondition: { prefix, suffix, ranges }, onChange }) => (
  <FormControl component='fieldset'>
    <FormLabel component='legend'>{name}</FormLabel>
    <RadioGroup
      aria-label={name}
      name={name}
      value={value}
      onChange={onChange}
      row
    >
      {
        ranges.map(({ id, min, max }) => {
          const minLabel = min !== -1 ? `${prefix}${min}${suffix} ` : ''
          const maxLabel = max !== -1 ? ` ${prefix}${max}${suffix}` : ''
          return <FormControlLabel key={id} value={id.toString()} control={<Radio />} label={`${minLabel}〜${maxLabel}`} />
        })
      }
      <FormControlLabel value='' control={<Radio />} label='指定なし' />
    </RadioGroup>
  </FormControl>
)
