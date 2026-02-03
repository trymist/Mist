import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Trash2, Plus, Pencil, X, Check, ExternalLink, ChevronDown, ChevronUp, AlertTriangle, RefreshCw } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { useDomains } from "@/hooks";
import { applicationsService } from "@/services";
import { DNSValidation } from "./dns-validation";
import type { Domain } from "@/types";

interface DomainsProps {
  appId: number;
}

export const Domains = ({ appId }: DomainsProps) => {
  const { domains, loading, createDomain, updateDomain, deleteDomain, updateDomainInState } = useDomains({
    appId,
    autoFetch: true
  });

  const [newDomain, setNewDomain] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editDomain, setEditDomain] = useState("");
  const [showAddForm, setShowAddForm] = useState(false);
  const [expandedDomain, setExpandedDomain] = useState<number | null>(null);
  const [actionDialogOpen, setActionDialogOpen] = useState(false);
  const [isRestarting, setIsRestarting] = useState(false);
  const [restartError, setRestartError] = useState<string | null>(null);

  const showRestartDialog = () => {
    setActionDialogOpen(true);
  };

  const handleRestart = async () => {
    try {
      setIsRestarting(true);
      setRestartError(null);
      await applicationsService.recreateContainer(appId);
      toast.success("Container recreated successfully - domain changes applied");
      setActionDialogOpen(false);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : "Failed to recreate container";
      setRestartError(errorMessage);
    } finally {
      setIsRestarting(false);
    }
    const handleSkipRestart = () => {
      setActionDialogOpen(false);
      toast.info("Domain changes will take effect on next restart or redeployment");
    };

    const handleAdd = async (e: React.FormEvent) => {
      e.preventDefault();
      if (!newDomain.trim()) {
        toast.error("Domain is required");
        return;
      }

      const result = await createDomain(newDomain.trim());
      if (result) {
        setNewDomain("");
        setShowAddForm(false);
        // Show restart dialog since domains require container recreation
        showRestartDialog();
      }
    };

    const handleUpdate = async (id: number) => {
      if (!editDomain.trim()) {
        toast.error("Domain is required");
        return;
      }

      const result = await updateDomain(id, editDomain.trim());
      if (result) {
        setEditingId(null);
        // Show restart dialog since domains require container recreation
        showRestartDialog();
      }
    };

    const handleDelete = async (id: number) => {
      if (!confirm("Are you sure you want to delete this domain?")) {
        return;
      }
      const result = await deleteDomain(id);
      if (result) {
        // Show restart dialog since domains require container recreation
        showRestartDialog();
      }
    };

    const startEdit = (domain: Domain) => {
      setEditingId(domain.id);
      setEditDomain(domain.domain);
    };

    const cancelEdit = () => {
      setEditingId(null);
      setEditDomain("");
    };

    const getSslStatusColor = (status: string) => {
      switch (status) {
        case "active":
          return "default";
        case "pending":
          return "secondary";
        case "failed":
          return "destructive";
        default:
          return "outline";
      }
    };

    const getDnsStatusBadge = (domain: Domain) => {
      if (domain.dnsConfigured) {
        return (
          <Badge variant="default" className="text-xs">
            DNS Configured
          </Badge>
        );
      }
      return (
        <Badge variant="secondary" className="text-xs">
          DNS Pending
        </Badge>
      );
    };

    const toggleDomainExpansion = (domainId: number) => {
      setExpandedDomain(expandedDomain === domainId ? null : domainId);
    };

    if (loading) {
      return <div className="text-muted-foreground">Loading domains...</div>;
    }

    return (
      <Card>
        <CardHeader>
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div>
              <CardTitle>Domains</CardTitle>
              <CardDescription>
                Manage custom domains for your application. Changes will be applied on next deployment.
              </CardDescription>
            </div>
            {!showAddForm && (
              <Button onClick={() => setShowAddForm(true)} size="sm" className="w-full sm:w-auto">
                <Plus className="h-4 w-4 mr-2" />
                Add Domain
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {showAddForm && (
            <form onSubmit={handleAdd} className="space-y-4 p-4 border rounded-lg bg-muted/50">
              <div className="space-y-2">
                <Label htmlFor="new-domain">Domain</Label>
                <Input
                  id="new-domain"
                  placeholder="example.com"
                  value={newDomain}
                  onChange={(e) => setNewDomain(e.target.value)}
                  autoFocus
                />
                <p className="text-sm text-muted-foreground">
                  Enter the domain without http:// or https://
                </p>
              </div>
              <div className="flex flex-col sm:flex-row gap-2">
                <Button type="submit" size="sm" className="w-full sm:w-auto">
                  <Check className="h-4 w-4 mr-2" />
                  Add
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="w-full sm:w-auto"
                  onClick={() => {
                    setShowAddForm(false);
                    setNewDomain("");
                  }}
                >
                  <X className="h-4 w-4 mr-2" />
                  Cancel
                </Button>
              </div>
            </form>
          )}

          {domains.length === 0 && !showAddForm ? (
            <p className="text-muted-foreground text-center py-8">
              No domains added yet. Click "Add Domain" to get started.
            </p>
          ) : (
            <div className="space-y-2">
              {domains.map((domain) => (
                <div key={domain.id} className="border rounded-lg bg-card">
                  {editingId === domain.id ? (
                    <div className="p-4 space-y-4">
                      <div className="space-y-2">
                        <Label htmlFor={`edit-domain-${domain.id}`}>Domain</Label>
                        <Input
                          id={`edit-domain-${domain.id}`}
                          value={editDomain}
                          onChange={(e) => setEditDomain(e.target.value)}
                        />
                      </div>
                      <div className="flex flex-col sm:flex-row gap-2">
                        <Button size="sm" onClick={() => handleUpdate(domain.id)} className="w-full sm:w-auto">
                          <Check className="h-4 w-4 mr-2" />
                          Save
                        </Button>
                        <Button size="sm" variant="outline" onClick={cancelEdit} className="w-full sm:w-auto">
                          <X className="h-4 w-4 mr-2" />
                          Cancel
                        </Button>
                      </div>
                    </div>
                  ) : (
                    <>
                      <div className="p-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                        <div className="flex flex-col sm:flex-row sm:items-center gap-3 flex-1 min-w-0">
                          <a
                            href={`https://${domain.domain}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="font-mono font-semibold hover:underline flex items-center gap-1 break-all"
                          >
                            {domain.domain}
                            <ExternalLink className="h-3 w-3 flex-shrink-0" />
                          </a>
                          <div className="flex items-center gap-2 flex-wrap">
                            <Badge variant={getSslStatusColor(domain.sslStatus)}>
                              SSL: {domain.sslStatus}
                            </Badge>
                            {getDnsStatusBadge(domain)}
                          </div>
                        </div>
                        <div className="flex gap-2 self-end sm:self-auto">
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => toggleDomainExpansion(domain.id)}
                          >
                            {expandedDomain === domain.id ? (
                              <ChevronUp className="h-4 w-4" />
                            ) : (
                              <ChevronDown className="h-4 w-4" />
                            )}
                          </Button>
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => startEdit(domain)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => handleDelete(domain.id)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                      {expandedDomain === domain.id && (
                        <div className="px-4 pb-4">
                          <DNSValidation
                            domain={domain}
                            onVerified={(updatedDomain) => {
                              updateDomainInState(updatedDomain);
                            }}
                          />
                        </div>
                      )}
                    </>
                  )}
                </div>
              ))}
            </div>
          )}
        </CardContent>

        {/* Restart Dialog */}
        <Dialog open={actionDialogOpen} onOpenChange={(open) => {
          if (!isRestarting) {
            setActionDialogOpen(open);
            if (!open) {
              setRestartError(null);
            }
          }
        }}>
          <DialogContent onPointerDownOutside={(e) => isRestarting && e.preventDefault()}>
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2">
                <AlertTriangle className="h-5 w-5 text-yellow-500" />
                Restart Required
              </DialogTitle>
              <DialogDescription>
                {isRestarting
                  ? "Restarting container, please wait..."
                  : "Domain changes require restarting the container to take effect. Would you like to restart now?"}
              </DialogDescription>
            </DialogHeader>
            {restartError && (
              <div className="p-3 rounded-md bg-destructive/10 text-destructive text-sm break-words max-h-32 overflow-y-auto">
                Error: {restartError}
              </div>
            )}
            <DialogFooter className="flex gap-2 sm:justify-end">
              <Button variant="outline" onClick={handleSkipRestart} disabled={isRestarting}>
                Skip for Now
              </Button>
              <Button onClick={handleRestart} disabled={isRestarting} className="flex items-center gap-2">
                <RefreshCw className={`h-4 w-4 ${isRestarting ? 'animate-spin' : ''}`} />
                {isRestarting ? 'Restarting...' : 'Restart Now'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </Card>
    );
  }
}
