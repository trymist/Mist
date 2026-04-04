import React from 'react';
import { Loader2, RefreshCw } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useAuditLogs } from './hooks/useAuditLogs';
import { 
  AuditLogList, 
  AuditLogFilters, 
  PaginationControls 
} from './components';

export const AuditLogsPage = () => {
  const {
    logs,
    loading,
    error,
    total,
    page,
    totalPages,
    setPage,
    setResourceFilter,
    refetch,
  } = useAuditLogs();

  const [currentFilter, setCurrentFilter] = React.useState('all');

  const handleFilterChange = (value: string) => {
    setCurrentFilter(value);
    setResourceFilter(value);
    setPage(0); // Reset to first page on filter change
  };

  if (loading && logs.length === 0) {
    return (
      <div className="flex items-center justify-center h-96">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Audit Logs</h1>
          <p className="text-muted-foreground mt-1">
            Track all activities and changes in your system
          </p>
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={refetch}
          disabled={loading}
        >
          <RefreshCw className={`h-4 w-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
          Refresh
        </Button>
      </div>

      {error && (
        <Card className="border-red-500/50 bg-red-500/10">
          <CardContent className="pt-6">
            <p className="text-red-400">{error}</p>
          </CardContent>
        </Card>
      )}

      <AuditLogFilters
        resourceFilter={currentFilter}
        onResourceFilterChange={handleFilterChange}
        total={total}
      />

      <Card>
        <CardHeader>
          <CardTitle>Activity History</CardTitle>
          <CardDescription>
            Showing {logs.length} of {total} audit logs
          </CardDescription>
        </CardHeader>
        <CardContent>
          <AuditLogList logs={logs} loading={loading} />
          
          <PaginationControls
            page={page}
            totalPages={totalPages}
            onPageChange={setPage}
            loading={loading}
          />
        </CardContent>
      </Card>
    </div>
  );
};
