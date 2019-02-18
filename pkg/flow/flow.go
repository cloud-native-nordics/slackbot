package flow

import(
	"errors"
)

type FlowManager struct{
	OnGoing []*Flow
	Questions []*Question
}

func (fm *FlowManager) AddNew(ChannelId string, UserId string){
	fm.Add(&Flow{
		ChannelID: ChannelId,
		UserID: UserId,
		CurrentStep: 1,
		Questions: fm.Questions,
	})
}

func (fm *FlowManager) Add(f *Flow){
	fm.OnGoing = append(fm.OnGoing, f)
}

func (fm *FlowManager) Get(userID string) (*Flow, error){
	for _, flow := range fm.OnGoing{
		if(flow.UserID == userID){
			return flow, nil
		}
	}
	return nil, errors.New("Cant find flow")
}

func (fm *FlowManager) IsInFlow(userID string) bool{
	flow, _ := fm.Get(userID)
	return flow != nil
}


func (fm *FlowManager) Remove(userID string){
	for i, p := range fm.OnGoing {
		if p.UserID == userID {
			fm.OnGoing = append(fm.OnGoing[:i], fm.OnGoing[i+1:]...)
		}
	}
}

type Flow struct{
	UserID string
	ChannelID string
	CurrentStep int
	Questions []*Question
	WaitingAnswer bool
}

func (f *Flow) SetCurrentAnswer(answer string){
	currentQuestion, err := f.GetCurrentQuestion()
	if(err != nil){
		panic(err)
	}
	currentQuestion.Answer = answer

}

func (f *Flow) GetPreviousQuestion() (*Question, error) {
	if f.CurrentStep == 1 {
		return nil, nil
	}

	for _, ele := range f.Questions{
		if(ele.Order == f.CurrentStep-1){
			return ele, nil
		}
	}
	return nil, errors.New("Cant find current question")
}

func (f *Flow) GetCurrentQuestion() (*Question, error) {
	for _, ele := range f.Questions{
		if(ele.Order == f.CurrentStep){
			return ele, nil
		}
	}
	return nil, errors.New("Cant find current question")
}

func (f *Flow) NextQuestion() {
	f.CurrentStep++
}

func (f *Flow) IsDone() bool{
	if(f.CurrentStep > len(f.Questions)){
		return true
	}
	return false
}


type Question struct{
	Order int
	Question string
	Answer string
	IsFirst bool
}