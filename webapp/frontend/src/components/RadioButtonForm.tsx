import {
  FormControl,
  FormLabel,
  FormControlLabel,
  Radio,
  RadioGroup
} from '@material-ui/core'

import type { FC, ChangeEvent } from 'react'

interface Props {
  name: string
  value: string
  items: string[]
  onChange: (event: ChangeEvent<HTMLInputElement>, value: string) => void
}

export const RadioButtonForm: FC<Props> = ({ name, value, items, onChange }) => (
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
        items.map(item => (
          <FormControlLabel key={item} value={item} control={<Radio />} label={item} />
        ))
      }
      <FormControlLabel value='' control={<Radio />} label='指定なし' />
    </RadioGroup>
  </FormControl>
)
