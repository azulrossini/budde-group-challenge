import type { GlobalThemeOverrides } from 'naive-ui'

function cssVar(name: string): string {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

// buildNaiveThemeOverrides resolves colors.css's custom properties into a
// Naive UI theme object. Naive can't take var(...) references directly (it
// does its own color math), so this reads the computed values at startup —
// colors.css stays the only place color values are actually defined.
export function buildNaiveThemeOverrides(): GlobalThemeOverrides {
  const accent = cssVar('--color-accent')
  const accentHover = cssVar('--color-accent-hover')
  const warning = cssVar('--color-warning')
  const error = cssVar('--color-error')
  const success = cssVar('--color-success')
  const surface = cssVar('--color-surface')
  const surfaceRaised = cssVar('--color-surface-raised')
  const border = cssVar('--color-border')
  const text = cssVar('--color-text')
  const textMuted = cssVar('--color-text-muted')

  return {
    common: {
      primaryColor: accent,
      primaryColorHover: accentHover,
      primaryColorPressed: accent,
      primaryColorSuppl: accentHover,
      warningColor: warning,
      warningColorHover: warning,
      warningColorPressed: warning,
      warningColorSuppl: warning,
      errorColor: error,
      errorColorHover: error,
      errorColorPressed: error,
      errorColorSuppl: error,
      successColor: success,
      successColorHover: success,
      successColorPressed: success,
      successColorSuppl: success,
      textColorBase: text,
      textColor1: text,
      textColor2: textMuted,
      textColor3: textMuted,
      cardColor: surface,
      modalColor: surface,
      popoverColor: surfaceRaised,
      inputColor: surfaceRaised,
      tableColor: surface,
      borderColor: border,
      dividerColor: border,
    },
  }
}
