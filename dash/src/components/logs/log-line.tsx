import { memo } from 'react';
import { parseAnsi, type AnsiSegment } from '@/lib/ansi-parser';
import { cn } from '@/lib/utils';

export interface LogLineProps {
  line: string;
  index: number;
  showLineNumbers?: boolean;
  streamType?: 'stdout' | 'stderr';
  className?: string;
}

/**
 * Renders a single log line with ANSI color support and stream type indicators
 */
export const LogLine = memo(({
  line,
  index,
  showLineNumbers = true,
  streamType,
  className,
}: LogLineProps) => {
  const segments = parseAnsi(line);

  return (
    <div
      className={cn(
        'hover:bg-slate-900/50 px-2 py-0.5 rounded transition-colors flex items-start gap-2',
        className
      )}
    >
      {/* Line Number */}
      {showLineNumbers && (
        <span className="text-slate-500 select-none font-mono text-sm min-w-[3ch] text-right shrink-0 mt-[2px]">
          {String(index + 1).padStart(4, ' ')}
        </span>
      )}

      {/* Stream Type Indicator */}
      {streamType && (
        <span
          className={cn(
            'select-none font-mono text-xs px-1.5 py-1.5 rounded shrink-0 mt-[2px]',
            streamType === 'stderr'
              ? 'bg-red-500/20 text-red-400'
              : 'bg-blue-500/20 text-blue-400'
          )}
          title={streamType === 'stderr' ? 'Standard Error' : 'Standard Output'}
        >
          {streamType === 'stderr' ? 'ERR' : 'OUT'}
        </span>
      )}

      {/* Log Content with ANSI colors */}
      <div className="flex-1 whitespace-pre-wrap break-all font-mono text-sm leading-relaxed">
        {segments.map((segment: AnsiSegment, i: number) => {
          const style: React.CSSProperties = {};

          if (segment.styles.color) {
            style.color = segment.styles.color;
          }
          if (segment.styles.backgroundColor) {
            style.backgroundColor = segment.styles.backgroundColor;
          }
          if (segment.styles.bold) {
            style.fontWeight = 'bold';
          }
          if (segment.styles.italic) {
            style.fontStyle = 'italic';
          }
          if (segment.styles.underline) {
            style.textDecoration = 'underline';
          }
          if (segment.styles.dim) {
            style.opacity = 0.5;
          }

          return (
            <span className='text-base' key={i} style={style}>
              {segment.text}
            </span>
          );
        })}
      </div>
    </div>
  );
});

LogLine.displayName = 'LogLine';
