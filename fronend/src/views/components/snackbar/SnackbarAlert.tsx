import { Fragment, SyntheticEvent, useState } from 'react'

import Alert from '@mui/material/Alert'
import Button from '@mui/material/Button'
import Snackbar from '@mui/material/Snackbar'

import { useSettings } from 'src/@core/hooks/useSettings'

const SnackbarAlert = () => {
  const [open, setOpen] = useState<boolean>(false)
  const { settings } = useSettings()
  const { skin } = settings

  const handleClick = () => {
    setOpen(true)
  }

  const handleClose = (event?: Event | SyntheticEvent, reason?: string) => {
    if (reason === 'clickaway') {
      return
    }
    setOpen(false)
  }

  return (
    <Fragment>
      <Button variant='outlined' onClick={handleClick}>
        Open alert snackbar
      </Button>
      <Snackbar open={open} onClose={handleClose} autoHideDuration={3000}>
        <Alert
          variant='filled'
          severity='success'
          onClose={handleClose}
          sx={{ width: '100%' }}
          elevation={skin === 'bordered' ? 0 : 3}
        >
          This is a success message!
        </Alert>
      </Snackbar>
    </Fragment>
  )
}

export default SnackbarAlert
