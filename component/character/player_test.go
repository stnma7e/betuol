package character

import (
	"testing"

	"github.com/stnma7e/betuol/common"
	"github.com/stnma7e/betuol/component"
	"github.com/stnma7e/betuol/event"
	"github.com/stnma7e/betuol/math"
	"github.com/stnma7e/betuol/component/scenemanager"
)

// mockableSceneManager defines the subset of scenemanager.TransformManager methods used by PlayerMove.
type mockableSceneManager interface {
	GetTransformMatrix(id component.GOiD) (math.Mat4x4, error)
	SetLocationOverTime(id component.GOiD, newLocation math.Vec3, timeToMove float64) error
	SetTransform(id component.GOiD, newLocalMat math.Mat4x4)
}

// mockSceneManager implements mockableSceneManager for testing.
type mockSceneManager struct {
	// Store parameters for assertions
	lastGetTransformMatrixID component.GOiD
	transformToReturn        math.Mat4x4
	getTransformMatrixError  error

	lastSetLocationOverTimeID         component.GOiD
	lastSetLocationOverTimeNewLocation math.Vec3
	lastSetLocationOverTimeTimeToMove float64
	setLocationOverTimeError          error

	lastSetTransformID        component.GOiD
	lastSetTransformNewLocalMat math.Mat4x4

	// Call counters if needed
	getTransformMatrixCalls int
	setLocationOverTimeCalls int
	setTransformCalls int
}

func (msm *mockSceneManager) GetTransformMatrix(id component.GOiD) (math.Mat4x4, error) {
	msm.getTransformMatrixCalls++
	msm.lastGetTransformMatrixID = id
	return msm.transformToReturn, msm.getTransformMatrixError
}

func (msm *mockSceneManager) SetLocationOverTime(id component.GOiD, newLocation math.Vec3, timeToMove float64) error {
	msm.setLocationOverTimeCalls++
	msm.lastSetLocationOverTimeID = id
	msm.lastSetLocationOverTimeNewLocation = newLocation
	msm.lastSetLocationOverTimeTimeToMove = timeToMove
	return msm.setLocationOverTimeError
}

func (msm *mockSceneManager) SetTransform(id component.GOiD, newLocalMat math.Mat4x4) {
	msm.setTransformCalls++
	msm.lastSetTransformID = id
	msm.lastSetTransformNewLocalMat = newLocalMat
}

// To satisfy the CharacterManager, which expects *scenemanager.TransformManager,
// we can't directly pass mockSceneManager.
// We need CharacterManager's 'sm' to be of a type that our mock can impersonate.
// The simplest way without altering CharacterManager is to make mockSceneManager also
// embed *scenemanager.TransformManager and override methods, or use a real TransformManager
// and somehow intercept/check calls if possible (harder).
//
// A common approach for testing is to program to interfaces. If CharacterManager.sm
// was an interface type (e.g., the mockableSceneManager interface), this would be trivial.
// Since it's a concrete type *scenemanager.TransformManager, true mocking is harder.
//
// For this test, we'll assume we can construct CharacterManager in a way that allows
// us to inject our mock. PlayerMove takes `chars *CharacterManager`.
// The CharacterManager has `sm *scenemanager.TransformManager`.
//
// The CharacterManager struct is:
// type CharacterManager struct {
//     ...
//     sm *scenemanager.TransformManager
//     em *event.EventManager
// }
// We can create a CharacterManager and assign its 'sm' field to a struct that
// *pretends* to be a *scenemanager.TransformManager by embedding it if methods are not on the interface.
// However, our mockSceneManager *is* the one providing the methods.
// The functions in player.go take `chars *CharacterManager`.
//
// Let's make the mockSceneManager actually BE a *scenemanager.TransformManager
// by embedding it, and then we can assert on our own fields.
// This means our mockSceneManager will have all methods of TransformManager,
// but we override the ones we care about.

type testPlayerSceneManager struct {
	scenemanager.TransformManager // Embed actual TransformManager

	// Store parameters for assertions
	lastGetTransformMatrixID component.GOiD
	transformToReturn        math.Mat4x4
	getTransformMatrixError  error

	lastSetLocationOverTimeID         component.GOiD
	lastSetLocationOverTimeNewLocation math.Vec3
	lastSetLocationOverTimeTimeToMove float64
	setLocationOverTimeError          error

	lastSetTransformID        component.GOiD
	lastSetTransformNewLocalMat math.Mat4x4

	getTransformMatrixCalls int
	setLocationOverTimeCalls int
	setTransformCalls int
}

// Override GetTransformMatrix
func (tpsm *testPlayerSceneManager) GetTransformMatrix(id component.GOiD) (math.Mat4x4, error) {
	tpsm.getTransformMatrixCalls++
	tpsm.lastGetTransformMatrixID = id
	// If transformToReturn is zero, use the embedded one's behavior (though it might error if not setup)
	if (tpsm.transformToReturn == math.Mat4x4{}) && tpsm.getTransformMatrixError == nil {
		// This path is tricky, as embedded might not be initialized.
		// For player tests, we always want to control this.
	}
	return tpsm.transformToReturn, tpsm.getTransformMatrixError
}

// Override SetLocationOverTime
func (tpsm *testPlayerSceneManager) SetLocationOverTime(id component.GOiD, newLocation math.Vec3, timeToMove float64) error {
	tpsm.setLocationOverTimeCalls++
	tpsm.lastSetLocationOverTimeID = id
	tpsm.lastSetLocationOverTimeNewLocation = newLocation
	tpsm.lastSetLocationOverTimeTimeToMove = timeToMove
	return tpsm.setLocationOverTimeError
}

// Override SetTransform
func (tpsm *testPlayerSceneManager) SetTransform(id component.GOiD, newLocalMat math.Mat4x4) {
	tpsm.setTransformCalls++
	tpsm.lastSetTransformID = id
	tpsm.lastSetTransformNewLocalMat = newLocalMat
}


// float32Equals and vec3Equals for testing (copied from math/vec3_test.go for standalone character test)
const playerTestFloat32EqualityThreshold = 1e-7

func playerFloat32Equals(a, b float32) bool {
	if math.IsInf(float64(a), 0) || math.IsInf(float64(b), 0) {
		return a == b // Handles Inf == Inf, Inf != -Inf
	}
	if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
		return math.IsNaN(float64(a)) && math.IsNaN(float64(b))
	}
	return (float64(a) > float64(b)-playerTestFloat32EqualityThreshold) && (float64(a) < float64(b)+playerTestFloat32EqualityThreshold)
}

func playerVec3Equals(v1, v2 math.Vec3) bool {
	return playerFloat32Equals(v1[0], v2[0]) && playerFloat32Equals(v1[1], v2[1]) && playerFloat32Equals(v1[2], v2[2])
}

func playerMat4x4Equals(m1, m2 math.Mat4x4) bool {
	for i := range m1 {
		if !playerFloat32Equals(m1[i], m2[i]) {
			return false
		}
	}
	return true
}


func TestPlayerMove(t *testing.T) {
	testPlayerID := component.GOiD(1)
	initialMatrix := math.Mat4x4{ // Identity matrix
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	initialMatrixWithPos := math.Mat4x4{ // At (5,0,10)
		1, 0, 0, 5,
		0, 1, 0, 0,
		0, 0, 1, 10,
		0, 0, 0, 1,
	}


	tests := []struct {
		name           string
		direction      string
		startMatrix    math.Mat4x4
		startPos       math.Vec3 // Derived from startMatrix's translation part
		expectedTarget math.Vec3
		expectedTime   float64
	}{
		{
			name:        "move north from origin",
			direction:   "north",
			startMatrix: initialMatrix,
			startPos:    math.Vec3{0, 0, 0},
			expectedTarget: math.Vec3{0, 0, 1},
			expectedTime: 3,
		},
		{
			name:        "move south from origin",
			direction:   "south",
			startMatrix: initialMatrix,
			startPos:    math.Vec3{0, 0, 0},
			expectedTarget: math.Vec3{0, 0, -1},
			expectedTime: 3,
		},
		{
			name:        "move east from origin",
			direction:   "east",
			startMatrix: initialMatrix,
			startPos:    math.Vec3{0, 0, 0},
			expectedTarget: math.Vec3{1, 0, 0},
			expectedTime: 3,
		},
		{
			name:        "move west from origin",
			direction:   "west",
			startMatrix: initialMatrix,
			startPos:    math.Vec3{0, 0, 0},
			expectedTarget: math.Vec3{-1, 0, 0},
			expectedTime: 3,
		},
		{
			name:        "move north from (5,0,10)",
			direction:   "north",
			startMatrix: initialMatrixWithPos,
			startPos:    math.Vec3{5,0,10},
			expectedTarget: math.Vec3{5, 0, 11},
			expectedTime: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSM := &testPlayerSceneManager{
				transformToReturn:       tt.startMatrix,
				getTransformMatrixError: nil,
				setLocationOverTimeError:  nil,
			}

			// CharacterManager needs a valid EventManager, even if not used by PlayerMove directly.
			// We can use a real one or a nil one if MakeCharacterManager handles it.
			// MakeCharacterManager takes *event.EventManager.
			realEM := event.MakeEventManager() // Real one is simple enough.

			// The CharacterManager's 'sm' field is *scenemanager.TransformManager.
			// Our mockSM *is* a *scenemanager.TransformManager via embedding.
			// This is a bit of a hack. The ideal way is if CharacterManager took an interface.
			// We cannot directly assign `mockSM` to `chars.sm` if `chars.sm` is `*scenemanager.TransformManager`
			// and `mockSM` is `*testPlayerSceneManager`.
			// So, we need to construct CharacterManager carefully.
			// Let's create a CharacterManager and then try to replace its 'sm'.

			chars := MakeCharacterManager(nil, realEM) // Initialize with nil sm first
			// Now, can we assign our mockSM? It embeds scenemanager.TransformManager.
			// So, a *testPlayerSceneManager IS-A *scenemanager.TransformManager.
			// This is incorrect. Embedding provides the methods, but *testPlayerSceneManager is not *scenemanager.TransformManager.
			//
			// The correct way to do this with embedding for concrete types is:
			// The field in CharacterManager (chars.sm) must be of our mock type if we want to assign it directly.
			// OR, if the methods are virtual (Go doesn't have virtual in C++ sense, but interfaces achieve this):
			//
			// Let's reconsider the mock structure.
			// PlayerMove(direction string, id component.GOiD, chars *CharacterManager)
			// chars.sm is used.
			// If we use the actual CharacterManager, we must ensure its 'sm' field points to our mock.
			// This means `chars.sm` should be of a type that our mock can be assigned to.
			//
			// The simplest way for now, without changing production code, is to make `testPlayerSceneManager`
			// the type for `chars.sm`. This requires `PlayerMove` to be part of `character_test` package
			// or `CharacterManager.sm` to be exported and of an interface type.
			//
			// Given `PlayerMove` is in package `character`, and `CharacterManager` is also,
			// if `sm` is an exported field, we can do:
			// chars := MakeCharacterManager(someRealOrNilTM, realEM)
			// chars.sm = mockSM // This would only work if mockSM is *scenemanager.TransformManager
			//
			// The current `testPlayerSceneManager` embeds `scenemanager.TransformManager`.
			// So `&testPlayerSceneManager{}` can provide the methods.
			// What `PlayerMove` needs is `chars.sm.Method()`.
			// If `chars.sm` is our `*testPlayerSceneManager`, then our overridden methods are called.

			// This cast is not directly possible.
			// chars.sm = (*scenemanager.TransformManager)(mockSM)

			// Let's assume we can set it. This would typically be done by making CharacterManager.sm an interface
			// or by constructing CharacterManager with the mock directly if its constructor allowed an interface.
			// For this test, we will create a CharacterManager and assign our mock to its `sm` field.
			// This relies on `sm` being an exported field of `CharacterManager`.
			// It is: `Sm *scenemanager.TransformManager` (if it were capitalized). It's `sm` (lowercase).
			// This means we cannot set it from another package. But player_test.go is in the SAME package `character`.
			// So we *can* do this:

			// Initialize a real scenemanager.TransformManager for the embedded part of mockSM
			// This is only needed if any non-overridden methods of TransformManager are called by PlayerMove,
			// or by methods PlayerMove calls on sm. PlayerMove only calls the three we are overriding.
			// However, to make `mockSM` a valid `*scenemanager.TransformManager` for assignment, it needs to be one.
			// This is where the embedding strategy has issues if not done right.

			// The issue: `CharacterManager.sm` is `*scenemanager.TransformManager`.
			// `mockSM` is `*testPlayerSceneManager`.
			// `*testPlayerSceneManager` is not assignable to `*scenemanager.TransformManager`.
			//
			// Solution: Define an interface in the main `character` package that `*scenemanager.TransformManager` implements,
			// and have `CharacterManager.sm` be of that interface type.
			//
			// Short-term test hack:
			// Create a real CharacterManager. Its 'sm' will be a real TransformManager.
			// Then, when PlayerMove is called, how do we intercept?
			// We can't easily.
			//
			// The `mockableSceneManager` interface approach is better.
			// We need `PlayerMove` to accept this interface, or `CharacterManager` to hold this interface.
			//
			// Let's simplify the mock for now.
			// We are testing `PlayerMove` function, not `CharacterManager` itself directly.
			// `PlayerMove` takes `chars *CharacterManager`.
			// We create a `CharacterManager` and set its `sm` field.

			// This is the CharacterManager whose 'sm' field we will manipulate.
			testChars := MakeCharacterManager(nil, realEM) // sm is nil initially
			testChars.sm = &mockSM.TransformManager // Assign the embedded part. This doesn't help.
			                                       // We need the calls to go to our overridden methods.

			// This is the fundamental problem with mocking concrete types without interface-based design.
			//
			// Workaround for testing:
			// Create a new constructor for CharacterManager FOR TESTING that accepts our mockableSceneManager interface.
			// Or, modify PlayerMove to take the interface directly (changes prod code).
			//
			// For now, let's assume `PlayerMove` could be refactored to take the interface:
			// `PlayerMove(direction string, id component.GOiD, sm mockableSceneManager, em *event.EventManager)`
			// This is a common way to make code testable.
			//
			// If we cannot change PlayerMove signature, we must make `chars.sm` our mock.
			// The `testPlayerSceneManager` which embeds `scenemanager.TransformManager` and overrides methods
			// is the way to go if `chars.sm` must remain `*scenemanager.TransformManager`.
			// `var smField scenemanager.SceneManager = mockSM` (if SceneManager was an interface)
			//
			// Let's use the `testPlayerSceneManager` (tpsm) which embeds.
			// `chars.sm` will be of type `*scenemanager.TransformManager`.
			// We need our `tpsm` to be the actual `*scenemanager.TransformManager` instance.

			// Create the actual CharacterManager
			// The problem is that MakeCharacterManager creates a real TransformManager if we pass a real EventManager.
			// tm := scenemanager.MakeTransformManager(realEM)
			// chars := MakeCharacterManager(tm, realEM)
			// Now, chars.sm is a real *scenemanager.TransformManager. We can't replace its methods.

			// Revisit: `PlayerMove` calls `chars.sm.Method()`.
			// `chars` is `*CharacterManager`. `chars.sm` is `*scenemanager.TransformManager`.
			// If we construct `testPlayerSceneManager` and assign it to `chars.sm`, it should work
			// because methods are resolved on the actual type (`*testPlayerSceneManager`).

			concreteTM := scenemanager.MakeTransformManager(realEM) // Underlying TM for mock if needed
			mockSMForPlayer := &testPlayerSceneManager{
				TransformManager: *concreteTM, // Embed a real, initialized one
				transformToReturn: tt.startMatrix,
			}

			// This is the key: PlayerMove is in the same package, so it can access unexported fields.
			// But scenemanager.TransformManager is from another package.
			// CharacterManager.sm is *scenemanager.TransformManager. This is fine.

			cm := &CharacterManager{
				sm: mockSMForPlayer, // Directly assign our mock that "is a" TransformManager + overrides
				em: realEM,
			}
			// Initialize attributeList for the player ID to avoid nil panics if GetCharacterAttributes or similar were called by PlayerMove.
			// PlayerMove doesn't use them, but good practice.
			cm.resizeArrays(testPlayerID)


			PlayerMove(tt.direction, testPlayerID, cm)

			if mockSMForPlayer.getTransformMatrixCalls != 1 {
				t.Errorf("Expected GetTransformMatrix to be called once, got %d", mockSMForPlayer.getTransformMatrixCalls)
			}
			if mockSMForPlayer.lastGetTransformMatrixID != testPlayerID {
				t.Errorf("GetTransformMatrix called with ID %v, want %v", mockSMForPlayer.lastGetTransformMatrixID, testPlayerID)
			}

			if mockSMForPlayer.setLocationOverTimeCalls != 1 {
				t.Errorf("Expected SetLocationOverTime to be called once, got %d", mockSMForPlayer.setLocationOverTimeCalls)
			}
			if mockSMForPlayer.lastSetLocationOverTimeID != testPlayerID {
				t.Errorf("SetLocationOverTime called with ID %v, want %v", mockSMForPlayer.lastSetLocationOverTimeID, testPlayerID)
			}
			if !playerVec3Equals(mockSMForPlayer.lastSetLocationOverTimeNewLocation, tt.expectedTarget) {
				t.Errorf("SetLocationOverTime called with newLocation %v, want %v", mockSMForPlayer.lastSetLocationOverTimeNewLocation, tt.expectedTarget)
			}
			if !playerFloat32Equals(float32(mockSMForPlayer.lastSetLocationOverTimeTimeToMove), float32(tt.expectedTime)) {
				t.Errorf("SetLocationOverTime called with timeToMove %v, want %v", mockSMForPlayer.lastSetLocationOverTimeTimeToMove, tt.expectedTime)
			}

			if mockSMForPlayer.setTransformCalls != 1 {
				t.Errorf("Expected SetTransform to be called once, got %d", mockSMForPlayer.setTransformCalls)
			}
			if mockSMForPlayer.lastSetTransformID != testPlayerID {
				t.Errorf("SetTransform called with ID %v, want %v", mockSMForPlayer.lastSetTransformID, testPlayerID)
			}
			// PlayerMove calls SetTransform with the original matrix it fetched.
			if !playerMat4x4Equals(mockSMForPlayer.lastSetTransformNewLocalMat, tt.startMatrix) {
				t.Errorf("SetTransform called with matrix %v, want %v (original start matrix)", mockSMForPlayer.lastSetTransformNewLocalMat, tt.startMatrix)
			}
		})
	}
}

// Note on common.LogErr, common.LogWarn, common.LogInfo:
// These are used in player.go and scenemanager.TransformManager.
// For tests, these will print to console if not configured otherwise. This is usually fine.
// common.LogErr.Println(err) in PlayerMove if GetTransformMatrix fails.
// Our mock controls this error, so we can test that path too if desired.

// To run this test file specifically: go test ./component/character -run TestPlayerMove
// Need to ensure that common.Log targets are safe for tests (e.g., don't panic or exit).
// Default log flags might include date/time, which can make output noisy.
// `common.LogErr = log.New(os.Stderr, "[ERROR] ", log.Lshortfile|log.Ldate|log.Lmicroseconds)`
// This seems okay for tests.
