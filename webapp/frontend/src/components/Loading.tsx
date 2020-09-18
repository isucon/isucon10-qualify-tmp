import { Box, CircularProgress } from '@material-ui/core'

export const Loading = () => (
  <Box width={1} height='100vh' display='flex' justifyContent='center' alignItems='center'>
    <CircularProgress />
  </Box>
)
