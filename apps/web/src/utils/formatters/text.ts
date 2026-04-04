export function getInitials(name: string, maxLength: number = 2): string {
  if (!name) return '';
  
  const words = name.trim().split(/\s+/);
  const initials = words
    .slice(0, maxLength)
    .map(word => word.charAt(0).toUpperCase())
    .join('');
    
  return initials;
}
