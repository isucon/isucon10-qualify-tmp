import {
  FormControl,
  FormLabel,
  FormControlLabel,
  FormGroup,
  Checkbox
} from '@material-ui/core'

import type { FC, ChangeEvent } from 'react'

interface Props {
  name: string
  checkList: boolean[]
  selectList: string[]
  onChange: (event: ChangeEvent<HTMLInputElement>, checked: boolean, index: number) => void
}

export const CheckboxForm: FC<Props> = ({ name, checkList, selectList, onChange }) => (
  <FormControl component='fieldset'>
    <FormLabel component='legend'>{name}</FormLabel>
    <FormGroup row>
      {
        selectList.map((select, i) => (
          <FormControlLabel
            key={select}
            control={
              <Checkbox
                checked={checkList[i]}
                name={select}
                onChange={(event, checked) => {
                  onChange(event, checked, i)
                }}
              />
            }
            label={select}
          />
        ))
      }
    </FormGroup>
  </FormControl>
)
