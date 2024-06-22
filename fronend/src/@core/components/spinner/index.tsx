import { useTheme } from '@mui/material/styles'
import Box, { BoxProps } from '@mui/material/Box'

import { Blocks } from 'react-loader-spinner'

const FallbackSpinner = ({ sx }: { sx?: BoxProps['sx'] }) => {
  const theme = useTheme()

  return (
    <Box
      sx={{
        height: '100vh',
        display: 'flex',
        alignItems: 'center',
        flexDirection: 'column',
        justifyContent: 'center',
        ...sx
      }}
    >
      <Blocks
        visible={true}
        height='100'
        width='100'
        ariaLabel='blocks-loading'
        wrapperStyle={{}}
        wrapperClass='blocks-wrapper'
      />
    </Box>
  )
}

export default FallbackSpinner
