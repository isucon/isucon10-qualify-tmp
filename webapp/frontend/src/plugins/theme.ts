import { createMuiTheme } from '@material-ui/core/styles'
import { lightGreen } from '@material-ui/core/colors'

export default createMuiTheme(
  {
    palette: {
      primary: lightGreen
    },
    typography: {
      button: {
        textTransform: 'none'
      }
    },
    props: {
      MuiTextField: {
        variant: 'outlined'
      }
    }
  }
)
