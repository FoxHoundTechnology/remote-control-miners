import { Fragment, SyntheticEvent, useState } from 'react'

import Button from '@mui/material/Button'
import Snackbar from '@mui/material/Snackbar'
import IconButton from '@mui/material/IconButton'

import Icon from 'src/@core/components/icon'

const SnackbarSimple = () => {
  const [open, setOpen] = useState<boolean>(false)

  const handleClick = () => {
    setOpen(true)
  }

  const handleClose = (event: Event | SyntheticEvent, reason?: string) => {
    if (reason === 'clickaway') {
      return
    }
    setOpen(false)
  }

  return (
    <Fragment>
      <Button variant='outlined' onClick={handleClick}>
        Open simple snackbar
      </Button>
      <Snackbar
        open={open}
        onClose={handleClose}
        message='Note archived'
        autoHideDuration={3000}
        action={
          <Fragment>
            <Button size='small' onClick={handleClose}>
              Undo
            </Button>
            <IconButton size='small' aria-label='close' color='inherit' onClick={handleClose}>
              <Icon icon='mdi:close' fontSize={20} />
            </IconButton>
          </Fragment>
        }
      />
    </Fragment>
  )
}

export default SnackbarSimple
