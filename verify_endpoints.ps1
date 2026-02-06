# Esperar a que el servidor levante
Start-Sleep -Seconds 5

$baseUrl = "http://localhost:8080/api/v1"

# 1. Login to get Token
Write-Host "`n--- Logging In (Admin) ---"
$loginBody = @{
  email    = "test@example.com"
  password = "password123"
} | ConvertTo-Json

try {
  $loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $loginBody -ContentType "application/json" -ErrorAction Stop
  $token = $loginResponse.token
  Write-Host "Login Success!"
}
catch {
  Write-Host "Login Failed: $_"
  exit 1
}

$headers = @{ Authorization = "Bearer $token" }

# 2. Invite Staff
Write-Host "`n--- inviting Staff (Auditor) ---"
$inviteBody = @{
  email = "auditor_new@iso.com"
  role  = "Auditor_Interno"
} | ConvertTo-Json

try {
  $invResponse = Invoke-RestMethod -Uri "$baseUrl/organization/staff/invite" -Method Post -Headers $headers -Body $inviteBody -ContentType "application/json" -ErrorAction Stop
  Write-Host "Invite Success: $($invResponse | ConvertTo-Json)"
}
catch {
  Write-Host "Invite Failed: $_"
}

# 3. List Staff
Write-Host "`n--- Listing Staff ---"
try {
  $staff = Invoke-RestMethod -Uri "$baseUrl/organization/staff" -Method Get -Headers $headers -ErrorAction Stop
  Write-Host "Staff List: $($staff | ConvertTo-Json)"
  $newUserID = $staff[0].user_id # Assuming at least one (likely the new one or self if linked)
  # Actually list returns all, we pick the last one or filter? 
  # Let's hope the new one is there.
}
catch {
  Write-Host "List Staff Failed: $_"
}

# 4. Create Audit
Write-Host "`n--- Creating Audit ---"
$auditBody = @{ title = "ISO 9001:2015 Certification" } | ConvertTo-Json
try {
  $audit = Invoke-RestMethod -Uri "$baseUrl/audits" -Method Post -Headers $headers -Body $auditBody -ContentType "application/json" -ErrorAction Stop
  Write-Host "Audit Created: $($audit | ConvertTo-Json)"
  $auditID = $audit.id
}
catch {
  Write-Host "Create Audit Failed: $_"
  exit 1
}

# 5. Assign Staff to Audit
Write-Host "`n--- Assigning Staff to Audit ---"
# Need a UserID. If ListStaff returned data, use one. 
# We'll invite another one to be sure we have the ID if list json is complex
# Getting UserID from database via docker for reliability in this script
$newUserID_DB = docker exec iso_stack_db psql -U user_admin -d iso_audit_db -t -c "SELECT id FROM users WHERE email='auditor_new@iso.com';"
$newUserID_DB = $newUserID_DB.Trim()

if ($newUserID_DB) {
  $assignBody = @{
    user_id       = $newUserID_DB
    role_in_audit = "Auxiliar"
  } | ConvertTo-Json

  try {
    $assignResp = Invoke-RestMethod -Uri "$baseUrl/audits/$auditID/assign" -Method Post -Headers $headers -Body $assignBody -ContentType "application/json" -ErrorAction Stop
    Write-Host "Assign Success: $($assignResp | ConvertTo-Json)"
  }
  catch {
    Write-Host "Assign Failed: $_"
  }
}
else {
  Write-Host "Skipping Assign: User ID not found."
}

# 6. List My Audits
Write-Host "`n--- Listing My Audits ---"
try {
  $myAudits = Invoke-RestMethod -Uri "$baseUrl/projects/my-audits" -Method Get -Headers $headers -ErrorAction Stop
  Write-Host "My Audits: $($myAudits | ConvertTo-Json)"
}
catch {
  Write-Host "List My Audits Failed: $_"
}

# 7. Check Public Link (Need to fetch link from DB first or assume Assign returned it? Assign doesn't return it)
# We fetch from DB
Write-Host "`n--- Testing Public Link ---"
if ($newUserID_DB) {
  $tempLink = docker exec iso_stack_db psql -U user_admin -d iso_audit_db -t -c "SELECT temporary_link FROM audit_assignments WHERE user_id='$newUserID_DB' AND audit_id='$auditID';"
  $tempLink = $tempLink.Trim()
    
  if ($tempLink) {
    try {
      $publicAudit = Invoke-RestMethod -Uri "$baseUrl/public/access/$tempLink" -Method Get -ErrorAction Stop
      Write-Host "Public Audit Access Success: $($publicAudit | ConvertTo-Json)"
    }
    catch {
      Write-Host "Public Audit Access Failed: $_"
    }
  }
  else {
    Write-Host "No temp link found."
  }
}
