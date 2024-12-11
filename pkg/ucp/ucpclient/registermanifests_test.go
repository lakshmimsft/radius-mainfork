/*
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package ucpclient

/*
import (
	"context"
	"testing"

	"github.com/radius-project/radius/pkg/ucp/hostoptions"
	"github.com/radius-project/radius/pkg/ucp/server"
	"github.com/stretchr/testify/require"
)

func TestRegisterManifests(t *testing.T) {
	ctx := context.Background()

	// Setup temporary manifest directory and files
	manifestDir := "testdata"

	// Create sample manifest files in manifestDir
	// ...
	options := &server.Options{
		Config: &hostoptions.UCPConfig{
			Manifests: hostoptions.ManifestConfig{
				ManifestDirectory: manifestDir,
			},
		},
	}

	// Setup fake UCP client with fake servers
	fakeUCPClient, err := NewFakeUCPClient()
	require.NoError(t, err)


}
*/
/*
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	// armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	armpolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/radius-project/radius/pkg/to"
	"github.com/radius-project/radius/pkg/ucp/api/v20231001preview"
	ucpfake "github.com/radius-project/radius/pkg/ucp/api/v20231001preview/fake"
	"github.com/stretchr/testify/require"
)

type FakeUCPServer struct {
	ResourceProvidersServer *ucpfake.ResourceProvidersServer
	ResourceTypesServer     *ucpfake.ResourceTypesServer
	APIVersionsServer       *ucpfake.APIVersionsServer
	LocationsServer         *ucpfake.LocationsServer
}

/*
type MultiTransport struct {
	transports map[string]http.RoundTripper
}

func (mt *MultiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for prefix, transport := range mt.transports {
		if strings.HasPrefix(req.URL.Path, prefix) {
			return transport.RoundTrip(req)
		}
	}
	return nil, fmt.Errorf("no transport found for request path: %s", req.URL.Path)
}*/

func NewFakeUCPServer() (*FakeUCPServer, error) {

	// Create fake servers for each client
	resourceProvidersServer := &ucpfake.ResourceProvidersServer{
		BeginCreateOrUpdate: func(
			ctx context.Context,
			planeName string,
			resourceProviderName string,
			resource v20231001preview.ResourceProviderResource,
			options *v20231001preview.ResourceProvidersClientBeginCreateOrUpdateOptions,
		) (resp azfake.PollerResponder[v20231001preview.ResourceProvidersClientCreateOrUpdateResponse], errResp azfake.ErrorResponder) {
			// Simulate successful creation
			result := v20231001preview.ResourceProvidersClientCreateOrUpdateResponse{
				ResourceProviderResource: resource,
			}
			resp.AddNonTerminalResponse(http.StatusCreated, nil)
			resp.SetTerminalResponse(http.StatusOK, result, nil)

			return
		},
		Get: func(
			ctx context.Context,
			planeName string,
			resourceProviderName string,
			options *v20231001preview.ResourceProvidersClientGetOptions,
		) (resp azfake.Responder[v20231001preview.ResourceProvidersClientGetResponse], errResp azfake.ErrorResponder) {
			response := v20231001preview.ResourceProvidersClientGetResponse{
				ResourceProviderResource: v20231001preview.ResourceProviderResource{
					Name: to.Ptr(resourceProviderName),
				},
			}
			resp.SetResponse(http.StatusOK, response, nil)
			return
		},
	}

	// Create other fake servers similarly
	resourceTypesServer := &ucpfake.ResourceTypesServer{
		BeginCreateOrUpdate: func(
			ctx context.Context,
			planeName string,
			resourceProviderName string,
			resourceTypeName string,
			resource v20231001preview.ResourceTypeResource,
			options *v20231001preview.ResourceTypesClientBeginCreateOrUpdateOptions,
		) (resp azfake.PollerResponder[v20231001preview.ResourceTypesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder) {
			result := v20231001preview.ResourceTypesClientCreateOrUpdateResponse{
				ResourceTypeResource: resource,
			}

			resp.AddNonTerminalResponse(http.StatusCreated, nil)
			resp.SetTerminalResponse(http.StatusOK, result, nil)

			return
		},
		Get: func(
			ctx context.Context,
			planeName string,
			resourceProviderName string,
			resourceTypeName string,
			options *v20231001preview.ResourceTypesClientGetOptions,
		) (resp azfake.Responder[v20231001preview.ResourceTypesClientGetResponse], errResp azfake.ErrorResponder) {
			response := v20231001preview.ResourceTypesClientGetResponse{
				ResourceTypeResource: v20231001preview.ResourceTypeResource{
					Name: to.Ptr(resourceTypeName),
				},
			}
			resp.SetResponse(http.StatusOK, response, nil)
			return
		},
	}

	apiVersionsServer := &ucpfake.APIVersionsServer{
		BeginCreateOrUpdate: func(
			ctx context.Context,
			planeName string,
			resourceProviderName string,
			resourceTypeName string,
			apiVersionName string, // Added missing parameter
			resource v20231001preview.APIVersionResource,
			options *v20231001preview.APIVersionsClientBeginCreateOrUpdateOptions,
		) (resp azfake.PollerResponder[v20231001preview.APIVersionsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder) {
			// Simulate successful creation
			result := v20231001preview.APIVersionsClientCreateOrUpdateResponse{
				APIVersionResource: resource,
			}
			resp.AddNonTerminalResponse(http.StatusCreated, nil)
			resp.SetTerminalResponse(http.StatusOK, result, nil)
			return
		},
	}

	locationsServer := &ucpfake.LocationsServer{
		BeginCreateOrUpdate: func(
			ctx context.Context,
			planeName string,
			resourceProviderName string,
			locationName string,
			resource v20231001preview.LocationResource,
			options *v20231001preview.LocationsClientBeginCreateOrUpdateOptions,
		) (resp azfake.PollerResponder[v20231001preview.LocationsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder) {
			// Simulate successful creation
			result := v20231001preview.LocationsClientCreateOrUpdateResponse{
				LocationResource: resource,
			}
			resp.AddNonTerminalResponse(http.StatusCreated, nil)
			resp.SetTerminalResponse(http.StatusOK, result, nil)

			return
		},
	}

	// Return the FakeUCPClient instance
	return &FakeUCPServer{
		ResourceProvidersServer: resourceProvidersServer,
		ResourceTypesServer:     resourceTypesServer,
		APIVersionsServer:       apiVersionsServer,
		LocationsServer:         locationsServer,
		/*UCPClient: ucpClient,

		ResourceProvidersClient: *resourceProvidersClient,
		ResourceTypesClient:     *resourceTypesClient,
		APIVersionsClient:       *apiVersionsClient,
		LocationsClient:         *locationsClient,
		*/
	}, nil
}

func NewTestUCPClient(server *FakeUCPServer) (*UCPClient, error) {

	// Create individual transports for each fake server
	resourceProvidersTransport := ucpfake.NewResourceProvidersServerTransport(server.ResourceProvidersServer)
	resourceTypesTransport := ucpfake.NewResourceTypesServerTransport(server.ResourceTypesServer)
	apiVersionsTransport := ucpfake.NewAPIVersionsServerTransport(server.APIVersionsServer)
	locationsTransport := ucpfake.NewLocationsServerTransport(server.LocationsServer)

	/*
	   // Create the ServerFactory
	   serverFactory := &ucpfake.ServerFactory{}

	   // Register your fake servers with the ServerFactory
	   serverFactory.ResourceProvidersServer = *resourceProvidersServer
	   serverFactory.ResourceTypesServer = *resourceTypesServer
	   serverFactory.APIVersionsServer = *apiVersionsServer
	   serverFactory.LocationsServer = *locationsServer
	   //locTransporter := ucpfake.ServerFactoryTransport(locationsServer)

	   serverFactoryTransport := ucpfake.NewServerFactoryTransport(serverFactory)
	*/

	//clientOptions := &policy.ClientOptions{Transport: serverFactoryTransport}

	// Assuming you have a UCPConnection or can use a fake one for testing
	//ucpConnection, err := sdk.NewDirectConnection("") // Replace with a valid implementation
	/*
	   // Create the server.Options instance
	   serverOptions := &server.Options{
	   	UCPConnection: ucpConnection,
	   }

	   httpClient:= http.Client{
	   	Transport: serverFactoryTransport,
	   }

	   ucpConnection.Client() = httpClient
	*/
	// Instantiate UCPClientImpl with the unified ClientOptions
	//ucpClient, err := NewUCPClient(ucpConnection)
	/*if err != nil {
		return nil, fmt.Errorf("failed to create UCP client: %w", err)
	}*/

	// Configure client options with respective transports
	resourceProvidersOptions := &armpolicy.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: resourceProvidersTransport,
		},
	}

	resourceTypesOptions := &armpolicy.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: resourceTypesTransport,
		},
	}

	apiVersionsOptions := &armpolicy.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: apiVersionsTransport,
		},
	}

	locationsOptions := &armpolicy.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Transport: locationsTransport,
		},
	}

	credential := &azfake.TokenCredential{}

	resourceProvidersClient, err := v20231001preview.NewResourceProvidersClient(credential, resourceProvidersOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create fake ResourceProvidersClient: %w", err)
	}

	resourceTypesClient, err := v20231001preview.NewResourceTypesClient(credential, resourceTypesOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create fake ResourceTypesClient: %w", err)
	}

	apiVersionsClient, err := v20231001preview.NewAPIVersionsClient(credential, apiVersionsOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create fake APIVersionsClient: %w", err)
	}

	locationsClient, err := v20231001preview.NewLocationsClient(credential, locationsOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create fake LocationsClient: %w", err)
	}

	return &UCPClient{
		ResourceProvidersClient: resourceProvidersClient,
		ResourceTypesClient:     resourceTypesClient,
		APIVersionsClient:       apiVersionsClient,
		LocationsClient:         locationsClient,
	}, nil
}

func TestRegisterManifests_Success(t *testing.T) {
	ctx := context.Background()
	expectedResourceProvider := "MyCompany2.CompanyName2"

	// Setup fake UCP server
	server, err := NewFakeUCPServer()
	require.NoError(t, err)

	// Setup UCP client with fake servers
	client, err := NewTestUCPClient(server)
	require.NoError(t, err)

	err = client.RegisterManifests(ctx, "testdata")
	require.NoError(t, err)

	// Verify resource provider was created
	rp, err := client.GetResourceProvider(ctx, planeName, expectedResourceProvider)
	require.NoError(t, err)
	require.Equal(t, to.Ptr(expectedResourceProvider), rp.Name)

}

func TestRegisterManifests_InvalidParameters(t *testing.T) {
	ctx := context.Background()

	// Setup fake UCP server
	server, err := NewFakeUCPServer()
	require.NoError(t, err)

	// Setup UCP client with fake servers
	client, err := NewTestUCPClient(server)
	require.NoError(t, err)

	// Pass invalid manifest directory
	err = client.RegisterManifests(ctx, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid manifest directory")
}
