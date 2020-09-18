import { useEffect } from 'react'
import Head from 'next/head'
import { ThemeProvider, makeStyles, createStyles } from '@material-ui/core/styles'
import CssBaseline from '@material-ui/core/CssBaseline'
import theme from '../plugins/theme'
import { NavBar } from '../components/NavBar'

import type { FC } from 'react'
import type { AppProps } from 'next/app'

import 'leaflet/dist/leaflet.css'

const useStyles = makeStyles(() =>
  createStyles({
    container: {
      width: '100vw',
      minHeight: '100vh',
      display: 'flex',
      flexDirection: 'column'
    }
  })
)

const MyApp: FC<AppProps> = props => {
  const { Component, pageProps } = props
  const classes = useStyles()

  useEffect(() => {
    // Remove the server-side injected CSS.
    const jssStyles = document.querySelector('#jss-server-side')
    if (jssStyles) {
      jssStyles.parentElement?.removeChild(jssStyles)
    }

    import('leaflet')
      .then(Leaflet => {
        delete (Leaflet.Icon.Default.prototype as any)._getIconUrl
        Leaflet.Icon.Default.mergeOptions({
          iconRetinaUrl: '/images/leaflet/marker-icon-2x.png',
          iconUrl: '/images/leaflet/marker-icon.png',
          shadowUrl: '/images/leaflet/marker-shadow.png'
        })
      })
      .catch(error => console.error('failed to set marker icons.', error))
  }, [])

  return (
    <>
      <Head>
        <title>isuumo</title>
        <meta name='viewport' content='minimum-scale=1, initial-scale=1, width=device-width' />
      </Head>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <div className={classes.container}>
          <NavBar />
          <Component {...pageProps} />
        </div>
      </ThemeProvider>
    </>
  )
}

export default MyApp
