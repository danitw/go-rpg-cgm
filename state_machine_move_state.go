package main

type MoveState struct {
	mCharacter  Character
	mMap        GameMap
	mEntity     *Entity
	mController *StateMachine
	// ^above common with WaitState
	mTileWidth       float64
	mMoveX, mMoveY   float64
	mPixelX, mPixelY float64
	mMoveSpeed       float64
	mTween           Tween
	mAnim            Animation
}

func MoveStateCreate(character Character, gMap GameMap) State {
	s := &MoveState{}
	s.mCharacter = character
	s.mMap = gMap
	s.mTileWidth = gMap.mTileWidth
	s.mEntity = character.mEntity
	s.mController = character.mController
	s.mMoveX = 0
	s.mMoveY = 0
	s.mTween = TweenCreate(0, 0, 0.5)
	s.mMoveSpeed = 0.4
	s.mAnim = AnimationCreate([]int{s.mEntity.startFrame}, true, 0.15)
	return s
}

//The StateMachine class requires each state to have
// four functions: Enter, Exit, Render and Update

func (s *MoveState) Enter(data Direction) {
	var frames []int
	if data.x == -1 {
		frames = s.mCharacter.mAnimLeft
	} else if data.x == 1 {
		frames = s.mCharacter.mAnimRight
	} else if data.y == -1 {
		frames = s.mCharacter.mAnimUp
	} else if data.y == 1 {
		frames = s.mCharacter.mAnimDown
	}
	s.mAnim.SetFrames(frames)

	//save Move X,Y value to used inside Update call
	s.mMoveX = data.x
	s.mMoveY = data.y
	s.mPixelX = s.mEntity.mTileX
	s.mPixelY = s.mEntity.mTileY
	s.mTween = TweenCreate(0, 1, s.mMoveSpeed)

	targetX, targetY := s.mEntity.mTileX+data.x, s.mEntity.mTileY+data.y
	if s.mMap.IsBlockingTile(int(targetX), int(targetY), 2) {
		s.mMoveX = 0
		s.mMoveY = 0
		s.mEntity.SetFrame(s.mAnim.Frame())
		s.mController.Change("wait", Direction{0, 0})
	}
}

func (s *MoveState) Exit() {
	if s.mMap.renderLayer == 2 {
		s.mEntity.TeleportAndDraw(s.mMap, s.mMap.canvas)
	}
}

func (s *MoveState) Render() {
	//pending
}

func (s *MoveState) Update(dt float64) {
	s.mAnim.Update(dt)
	s.mEntity.SetFrame(s.mAnim.Frame())

	s.mTween.Update(dt)
	value := s.mTween.Value()
	s.mEntity.mTileX = s.mPixelX + value*s.mMoveX
	s.mEntity.mTileY = s.mPixelY + value*s.mMoveY

	if s.mTween.IsFinished() {
		s.mController.Change("wait", Direction{0, 0})
	}
}