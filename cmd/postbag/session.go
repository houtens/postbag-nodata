package main

// const SessionExpiry int = 86400
// const SessionName string = "postbag-session"

// replace this with htmx
// func initLoginSession(c *gin.Context) {
// 	s := sessions.Default(c)
// 	if err := s.Save(); err != nil {
// 		fmt.Printf("error creating login session")
// 	}
// }

// func setLoginSession(c *gin.Context, data session.LoginSession) {
// 	s := sessions.Default(c)
// 	s.Set(SessionName, data)
// 	s.Options(sessions.Options{
// 		Path:   "/",
// 		MaxAge: SessionExpiry,
// 	})
// 	if err := s.Save(); err != nil {
// 		fmt.Println("unable to save login session")
// 	}
// }

// func getLoginSession(c *gin.Context) (session.LoginSession, error) {
// 	s := sessions.Default(c)
// 	data := s.Get(SessionName)
// 	if data != nil {
// 		// found a session, return it
// 		return data.(session.LoginSession), nil
// 	}
// 	if err := s.Save(); err != nil {
// 		return session.LoginSession{}, fmt.Errorf("error saving session: %w", err)
// 	}
// 	return session.LoginSession{}, session.ErrSessionNotFound
// }

// clear session on logout
// func clearSession(c *gin.Context) {
// 	s := sessions.Default(c)
// 	s.Clear()
// 	s.Options(sessions.Options{
// 		Path:   "/",
// 		MaxAge: -1,
// 	})
// 	_ = s.Save()
// }

// func setflash(c *gin.context, message string) {
// 	s := sessions.default(c)
// 	s.addflash(message)
// 	if err := s.save(); err != nil {
// 		fmt.printf("error saving flash message\n")
// 	}
// }

// type Flash struct {
// 	Message string
// 	Level   string
// }

// var Rx = regexp.MustCompile(`(info|success|warning|error):`)

// func parseFlash(msg string) Flash {
// 	// match warning level at the beginning of the message and split
// 	if Rx.MatchString(msg) {
// 		r := strings.Split(msg, ":")
// 		return Flash{
// 			Message: r[1],
// 			Level:   r[0],
// 		}
// 	}
// 	// no match, return as info
// 	return Flash{Message: msg, Level: "info"}
// }

// func getFlash(c *gin.Context) []Flash {
// 	s := sessions.Default(c)
// 	flashes := s.Flashes()
//
// 	flashesOutput := []Flash{}
//
// 	if len(flashes) != 0 {
// 		if err := s.Save(); err != nil {
// 			fmt.Printf("error reading flash messages")
// 		}
// 		for _, f := range flashes {
// 			flashesOutput = append(
// 				flashesOutput,
// 				parseFlash(f.(string)),
// 			)
// 		}
// 	}
// 	return flashesOutput
// }
