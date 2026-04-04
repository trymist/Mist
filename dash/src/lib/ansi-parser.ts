/**
 * ANSI Parser for Terminal Color Codes
 * Converts ANSI escape sequences to styled HTML spans
 */

export interface AnsiSegment {
  text: string;
  styles: {
    color?: string;
    backgroundColor?: string;
    bold?: boolean;
    italic?: boolean;
    underline?: boolean;
    dim?: boolean;
  };
}

// ANSI color mappings (standard 16 colors)
const ANSI_COLORS: Record<number, string> = {
  // Standard colors
  30: '#000000', // black
  31: '#ef4444', // red
  32: '#22c55e', // green
  33: '#eab308', // yellow
  34: '#3b82f6', // blue
  35: '#a855f7', // magenta
  36: '#06b6d4', // cyan
  37: '#f5f5f5', // white
  // Bright colors
  90: '#6b7280', // bright black (gray)
  91: '#f87171', // bright red
  92: '#4ade80', // bright green
  93: '#fbbf24', // bright yellow
  94: '#60a5fa', // bright blue
  95: '#c084fc', // bright magenta
  96: '#22d3ee', // bright cyan
  97: '#ffffff', // bright white
};

const ANSI_BG_COLORS: Record<number, string> = {
  // Background colors
  40: '#000000',
  41: '#ef4444',
  42: '#22c55e',
  43: '#eab308',
  44: '#3b82f6',
  45: '#a855f7',
  46: '#06b6d4',
  47: '#f5f5f5',
  // Bright background colors
  100: '#6b7280',
  101: '#f87171',
  102: '#4ade80',
  103: '#fbbf24',
  104: '#60a5fa',
  105: '#c084fc',
  106: '#22d3ee',
  107: '#ffffff',
};

interface AnsiState {
  color?: string;
  backgroundColor?: string;
  bold: boolean;
  italic: boolean;
  underline: boolean;
  dim: boolean;
}

/**
 * Parse ANSI escape sequences from a string and return styled segments
 */
export function parseAnsi(text: string): AnsiSegment[] {
  // First, strip out non-SGR ANSI escape sequences (cursor movement, clear screen, etc.)
  // Keep only SGR sequences (those ending in 'm')
  text = text.replace(/\x1b\[[\d;]*[A-HJKSTfhilmnsu]/g, '');
  
  const segments: AnsiSegment[] = [];
  const state: AnsiState = {
    bold: false,
    italic: false,
    underline: false,
    dim: false,
  };

  // Match ANSI escape sequences: \x1b[...m or \u001b[...m
  const ansiRegex = /\x1b\[([0-9;]*)m/g;
  let lastIndex = 0;
  let match: RegExpExecArray | null;

  while ((match = ansiRegex.exec(text)) !== null) {
    // Add text before this escape sequence
    if (match.index > lastIndex) {
      const textSegment = text.substring(lastIndex, match.index);
      if (textSegment) {
        segments.push({
          text: textSegment,
          styles: { ...state },
        });
      }
    }

    // Parse the escape codes
    const codes = match[1].split(';').map((code) => parseInt(code, 10));
    for (const code of codes) {
      if (isNaN(code)) continue;

      switch (code) {
        case 0: // Reset
          state.color = undefined;
          state.backgroundColor = undefined;
          state.bold = false;
          state.italic = false;
          state.underline = false;
          state.dim = false;
          break;
        case 1: // Bold
          state.bold = true;
          break;
        case 2: // Dim
          state.dim = true;
          break;
        case 3: // Italic
          state.italic = true;
          break;
        case 4: // Underline
          state.underline = true;
          break;
        case 22: // Normal intensity (not bold or dim)
          state.bold = false;
          state.dim = false;
          break;
        case 23: // Not italic
          state.italic = false;
          break;
        case 24: // Not underlined
          state.underline = false;
          break;
        case 39: // Default foreground color
          state.color = undefined;
          break;
        case 49: // Default background color
          state.backgroundColor = undefined;
          break;
        default:
          // Foreground colors (30-37, 90-97)
          if (ANSI_COLORS[code]) {
            state.color = ANSI_COLORS[code];
          }
          // Background colors (40-47, 100-107)
          else if (ANSI_BG_COLORS[code]) {
            state.backgroundColor = ANSI_BG_COLORS[code];
          }
          break;
      }
    }

    lastIndex = match.index + match[0].length;
  }

  // Add remaining text
  if (lastIndex < text.length) {
    const textSegment = text.substring(lastIndex);
    if (textSegment) {
      segments.push({
        text: textSegment,
        styles: { ...state },
      });
    }
  }

  // If no segments were created, return the whole text as one segment
  if (segments.length === 0) {
    segments.push({
      text,
      styles: {},
    });
  }

  return segments;
}

/**
 * Strip ANSI escape sequences from text
 */
export function stripAnsi(text: string): string {
  // Remove all ANSI escape sequences including SGR (m), cursor movement, etc.
  return text.replace(/\x1b\[[0-9;]*[A-HJKSTfhilmnsu]?/g, '');
}

/**
 * Detect if text contains ANSI escape sequences
 */
export function hasAnsi(text: string): boolean {
  return /\x1b\[[0-9;]*[A-HJKSTfhilmnsu]?/.test(text);
}
