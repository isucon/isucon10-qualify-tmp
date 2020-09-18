import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/router'
import { makeStyles, createStyles } from '@material-ui/core/styles'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import Menu from '@material-ui/core/Menu'
import MenuItem from '@material-ui/core/MenuItem'
import IconButton from '@material-ui/core/IconButton'
import MenuIcon from '@material-ui/icons/Menu'
import EventSeatIcon from '@material-ui/icons/EventSeat'
import HouseIcon from '@material-ui/icons/House'
import TouchAppIcon from '@material-ui/icons/TouchApp'

import type { FC } from 'react'

const useStyles = makeStyles(theme =>
  createStyles({
    menuButton: {
      marginRight: theme.spacing(2)
    },
    logo: {
      height: 48,
      cursor: 'pointer'
    }
  })
)

export const NavBar: FC = () => {
  const classes = useStyles()
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null)
  const router = useRouter()

  return (
    <AppBar position='relative'>
      <Toolbar>
        <IconButton
          edge='start'
          aria-haspopup='true'
          className={classes.menuButton}
          onClick={event => { setAnchorEl(event.currentTarget) }}
        >
          <MenuIcon />
        </IconButton>
        <Menu
          anchorEl={anchorEl}
          keepMounted
          open={Boolean(anchorEl)}
          onClose={() => { setAnchorEl(null) }}
        >
          <MenuItem onClick={async () => { await router.push('/chair/search') }}>
            <EventSeatIcon /> イス検索
          </MenuItem>
          <MenuItem onClick={async () => { await router.push('/estate/search') }}>
            <HouseIcon /> 物件検索
          </MenuItem>
          <MenuItem onClick={async () => { await router.push('/estate/nazotte') }}>
            <TouchAppIcon /> なぞって検索
          </MenuItem>
        </Menu>
        <Link href='/'>
          <img className={classes.logo} src='/images/logo.png' />
        </Link>
      </Toolbar>
    </AppBar>
  )
}
