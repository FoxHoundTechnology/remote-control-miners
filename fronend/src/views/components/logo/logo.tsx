import { PaletteMode } from '@mui/material'
import { styled, useTheme } from '@mui/material/styles'

const CLIENT_NAME = 'bwd'

const HeaderLogo = styled('img')(({ theme }) => ({
  height: 48,
  marginLeft: 40,
  marginTop: 20,
  zIndex: 9
}))

const FormLogo = styled('img')(({ theme }) => ({
  height: 98,
  width: 'auto',
  marginTop: 16,
  marginBottom: 16,
  zIndex: 9
}))

interface LogoProps {
  mode: PaletteMode
}

export const VerticalNavHeaderLogo: React.FC<LogoProps> = ({ mode }) => {
  return <HeaderLogo src={`/images/logos/ctg_${mode === 'dark' ? 'white' : 'black'}_logo.svg`} alt='logo' />
}

export const LoginScreenLogo: React.FC<LogoProps> = ({ mode }) => {
  return <FormLogo src={`/images/logos/ctg_${mode === 'dark' ? 'white' : 'black'}_logo.svg`} alt='logo' />
}
