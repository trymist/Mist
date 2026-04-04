import { useState, useEffect, useCallback } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { 
  CheckCircle2, 
  XCircle, 
  AlertCircle, 
  RefreshCw, 
  Copy, 
  Check 
} from "lucide-react";
import { toast } from "sonner";
import { applicationsService } from "@/services";
import type { Domain, DNSInstructions, DNSVerificationResponse } from "@/types";

interface DNSValidationProps {
  domain: Domain;
  onVerified?: (updatedDomain: Domain) => void;
}

export const DNSValidation = ({ domain, onVerified }: DNSValidationProps) => {
  const [instructions, setInstructions] = useState<DNSInstructions | null>(null);
  const [verifying, setVerifying] = useState(false);
  const [verificationResult, setVerificationResult] = useState<DNSVerificationResponse | null>(null);
  const [copiedRecord, setCopiedRecord] = useState<string | null>(null);
  const [currentDomain, setCurrentDomain] = useState<Domain>(domain);

  useEffect(() => {
    setCurrentDomain(domain);
  }, [domain]);

  const loadInstructions = useCallback(async () => {
    try {
      const data = await applicationsService.getDNSInstructions(currentDomain.id);
      setInstructions(data);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to load DNS instructions');
    }
  }, [currentDomain.id]);

  useEffect(() => {
    loadInstructions();
  }, [loadInstructions]);

  const handleVerify = async () => {
    setVerifying(true);
    try {
      const result = await applicationsService.verifyDomainDNS(currentDomain.id);
      setVerificationResult(result);
      
      if (result.valid) {
        toast.success('DNS configuration verified successfully!');
        // Update the current domain with the verified status
        setCurrentDomain(result.domain);
        onVerified?.(result.domain);
      } else {
        toast.error(result.error || 'DNS verification failed');
        // Update the domain with the error status
        setCurrentDomain(result.domain);
        onVerified?.(result.domain);
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to verify DNS');
    } finally {
      setVerifying(false);
    }
  };

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    setCopiedRecord(label);
    toast.success(`${label} copied to clipboard`);
    setTimeout(() => setCopiedRecord(null), 2000);
  };

  const getDnsStatusBadge = () => {
    if (currentDomain.dnsConfigured) {
      return (
        <Badge variant="default" className="gap-1">
          <CheckCircle2 className="h-3 w-3" />
          Verified
        </Badge>
      );
    }
    
    if (currentDomain.lastDnsCheck && !currentDomain.dnsConfigured) {
      return (
        <Badge variant="destructive" className="gap-1">
          <XCircle className="h-3 w-3" />
          Not Configured
        </Badge>
      );
    }

    return (
      <Badge variant="secondary" className="gap-1">
        <AlertCircle className="h-3 w-3" />
        Pending Verification
      </Badge>
    );
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              DNS Configuration
              {getDnsStatusBadge()}
            </CardTitle>
            <CardDescription>
              Configure your domain's DNS records to point to this server
            </CardDescription>
          </div>
          <Button 
            onClick={handleVerify} 
            disabled={verifying}
            size="sm"
          >
            {verifying ? (
              <>
                <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                Verifying...
              </>
            ) : (
              <>
                <RefreshCw className="h-4 w-4 mr-2" />
                Verify DNS
              </>
            )}
          </Button>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        {currentDomain.dnsConfigured && currentDomain.dnsVerifiedAt && (
          <Alert>
            <CheckCircle2 className="h-4 w-4" />
            <AlertDescription>
              DNS configured correctly and verified on{' '}
              {new Date(currentDomain.dnsVerifiedAt).toLocaleString()}
            </AlertDescription>
          </Alert>
        )}

        {currentDomain.dnsCheckError && !currentDomain.dnsConfigured && (
          <Alert variant="destructive">
            <XCircle className="h-4 w-4" />
            <AlertDescription>
              {currentDomain.dnsCheckError}
            </AlertDescription>
          </Alert>
        )}

        {verificationResult && !verificationResult.valid && verificationResult.error && (
          <Alert variant="destructive">
            <XCircle className="h-4 w-4" />
            <AlertDescription>
              {verificationResult.error}
            </AlertDescription>
          </Alert>
        )}

        {!currentDomain.dnsConfigured && (
          <div className="space-y-3">
            <div>
              <h4 className="font-semibold text-sm mb-2">Required DNS Records</h4>
              <p className="text-sm text-muted-foreground mb-3">
                Add the following DNS records to your domain's DNS settings:
              </p>
            </div>

            {instructions && (
              <div className="space-y-2">
                {instructions.records.map((record, index) => (
                  <div 
                    key={index} 
                    className="p-3 border rounded-lg bg-muted/50 space-y-2"
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Badge variant="outline">{record.type}</Badge>
                        <span className="font-mono text-sm font-semibold">
                          {record.name === '@' ? currentDomain.domain : `${record.name}.${currentDomain.domain}`}
                        </span>
                      </div>
                    </div>
                    <div className="flex items-center justify-between gap-2">
                      <code className="text-sm bg-background px-2 py-1 rounded flex-1">
                        {record.value}
                      </code>
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => copyToClipboard(record.value, `${record.type} record`)}
                      >
                        {copiedRecord === `${record.type} record` ? (
                          <Check className="h-4 w-4" />
                        ) : (
                          <Copy className="h-4 w-4" />
                        )}
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            )}

            <Alert>
              <AlertCircle className="h-4 w-4" />
              <AlertDescription className="text-sm">
                <strong>Note:</strong> DNS changes can take up to 48 hours to propagate, 
                though they typically complete within a few minutes. After updating your DNS 
                records, click "Verify DNS" to check the configuration.
              </AlertDescription>
            </Alert>
          </div>
        )}
      </CardContent>
    </Card>
  );
};
