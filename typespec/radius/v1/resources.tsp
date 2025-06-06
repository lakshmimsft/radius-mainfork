/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
    
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import "@azure-tools/typespec-autorest";
import "@azure-tools/typespec-azure-core";
import "@azure-tools/typespec-azure-resource-manager";
import "@typespec/openapi";

using OpenAPI;
using Azure.Core;
using Azure.ResourceManager;

@doc("The resource name")
@maxLength(63)
@pattern("^[A-Za-z]([-A-Za-z0-9]*[A-Za-z0-9])?$")
scalar ResourceNameString extends string;

@doc("Base properties of a Environment-scoped resource")
model EnvironmentScopedResource {
  @doc("Fully qualified resource ID for the environment that the portable resource is linked to")
  environment: string;

  @doc("Fully qualified resource ID for the application that the portable resource is consumed by (if applicable)")
  application?: string;

  @doc("The status of the asynchronous operation.")
  @visibility("read")
  provisioningState?: ProvisioningState;

  @doc("Status of a resource.")
  @visibility("read")
  status?: ResourceStatus;
}

@doc("Base properties of a Application-scoped resource")
model ApplicationScopedResource {
  @doc("Fully qualified resource ID for the environment that the application is linked to")
  environment?: string;

  @doc("Fully qualified resource ID for the application")
  application: string;

  @doc("The status of the asynchronous operation.")
  @visibility("read")
  provisioningState?: ProvisioningState;

  @doc("Status of a resource.")
  @visibility("read")
  status?: ResourceStatus;
}

@doc("Base properties of a Global-scoped resource")
model GlobalScopedResource {
  @doc("Fully qualified resource ID for the environment that the application is linked to")
  environment?: string;

  @doc("Fully qualified resource ID for the application")
  application?: string;

  @doc("The status of the asynchronous operation.")
  @visibility("read")
  provisioningState?: ProvisioningState;

  @doc("Status of a resource.")
  @visibility("read")
  status?: ResourceStatus;
}

@doc("IdentitySettings is the external identity setting.")
model IdentitySettings {
  @doc("kind of identity setting")
  kind: IdentitySettingKind;

  @doc("The URI for your compute platform's OIDC issuer")
  oidcIssuer?: string;

  @doc("The resource ID of the provisioned identity")
  resource?: string;

  @doc("The list of user assigned managed identities")
  managedIdentity?: string[];
}

@doc("IdentitySettingKind is the kind of supported external identity setting")
enum IdentitySettingKind {
  @doc("undefined identity")
  undefined: "undefined",

  @doc("azure ad workload identity")
  `azure.com.workload`: "azure.com.workload",

  @doc("User assigned managed identity")
  userAssigned: "userAssigned",

  @doc("System assigned managed identity")
  systemAssigned: "systemAssigned",

  @doc("System assigned and user assigned managed identity")
  systemAssignedUserAssigned: "systemAssignedUserAssigned",
}

@discriminator("kind")
@doc("Represents backing compute resource")
model EnvironmentCompute {
  @doc("The resource id of the compute resource for application environment.")
  resourceId?: string;

  @doc("Configuration for supported external identity providers")
  identity?: IdentitySettings;
}

@doc("The Kubernetes compute configuration")
model KubernetesCompute extends EnvironmentCompute {
  @doc("The Kubernetes compute kind")
  kind: "kubernetes";

  @doc("The namespace to use for the environment.")
  `namespace`: string;
}

@doc("The Azure container instance compute configuration")
model AzureContainerInstanceCompute extends EnvironmentCompute {
  @doc("The Azure container instance compute kind")
  kind: "aci";

  @doc("The resource group to use for the environment.")
  resourceGroup?: string;
}

@doc("Recipe status at deployment time for a resource.")
model RecipeStatus {
  @doc("TemplateKind is the kind of the recipe template used by the portable resource upon deployment.")
  templateKind: string;

  @doc("TemplatePath is the path of the recipe consumed by the portable resource upon deployment.")
  templatePath: string;

  @doc("TemplateVersion is the version number of the template.")
  templateVersion?: string;
}

@doc("Status of a resource.")
model ResourceStatus {
  @doc("The compute resource associated with the resource.")
  compute?: EnvironmentCompute;

  @doc("The recipe data at the time of deployment")
  @visibility("read")
  recipe?: RecipeStatus;

  @doc("Properties of an output resource")
  @extension("x-ms-identifiers", [])
  outputResources?: OutputResource[];
}

@doc("Properties of an output resource.")
model OutputResource {
  @doc("The logical identifier scoped to the owning Radius resource. This is only needed or used when a resource has a dependency relationship. LocalIDs do not have any particular format or meaning beyond being compared to determine dependency relationships.")
  localId?: string;

  @doc("The UCP resource ID of the underlying resource.")
  id?: string;

  @doc("Determines whether Radius manages the lifecycle of the underlying resource.")
  radiusManaged?: boolean;
}

@doc("Describes a reference to an existing resource")
model ResourceReference {
  @doc("Resource id of an existing resource")
  id: string;
}

@doc("Provisioning state of the resource at the time the operation was called")
@lroStatus
enum ProvisioningState {
  @doc("The resource is being created")
  Creating,

  @doc("The resource is being updated")
  Updating,

  @doc("The resource is being deleted")
  Deleting,

  @doc("The resource create request has been accepted")
  Accepted,

  @doc("The resource is being provisioned")
  Provisioning,

  @doc("The resource has been successfully provisioned")
  Succeeded,

  @doc("The resource provisioning has failed")
  Failed,

  @doc("The resource provisioning has been canceled")
  Canceled,
}

@doc("The base properties of Recipe")
model RecipeBaseProperties {
  @doc("The recipe used to automatically deploy underlying infrastructure for the resource")
  recipe?: Recipe;

  @doc("Specifies how the underlying service/resource is provisioned and managed.")
  resourceProvisioning?: ResourceProvisioning;
}

@doc("Specifies how the underlying service/resource is provisioned and managed. Available values are 'recipe', where Radius manages the lifecycle of the resource through a Recipe, and 'manual', where a user manages the resource and provides the values.")
enum ResourceProvisioning {
  @doc("The resource lifecycle will be managed by Radius")
  recipe,

  @doc("The resource lifecycle will be managed by the user")
  manual,
}

@doc("The recipe used to automatically deploy underlying infrastructure for a portable resource")
model Recipe {
  @doc("The name of the recipe within the environment to use")
  name: string;

  @doc("Key/value parameters to pass into the recipe at deployment")
  parameters?: {};
}
